package commands

import (
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/chen-keinan/beacon/internal/common"
	"github.com/chen-keinan/beacon/internal/logger"
	"github.com/chen-keinan/beacon/internal/models"
	"github.com/chen-keinan/beacon/internal/reports"
	"github.com/chen-keinan/beacon/internal/shell"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/kyokomi/emoji"
	"strconv"
	"strings"
)

var log = logger.GetLog()

//ValidateExprData expr data
type ValidateExprData struct {
	index     int
	resultArr []string
	atb       *models.AuditBench
	origSize  int
	Total     int
	Match     int
}

//NextValidExprData return the next recursive ValidExprData
func (ve ValidateExprData) NextValidExprData() ValidateExprData {
	return ValidateExprData{resultArr: ve.resultArr[1:ve.index], index: ve.index - 1, atb: ve.atb, origSize: ve.origSize}
}

// NewValidExprData return new instance of ValidExprData
func NewValidExprData(arr []string, at *models.AuditBench) ValidateExprData {
	return ValidateExprData{resultArr: arr, index: len(arr), atb: at, origSize: len(arr)}
}

//K8sAudit k8s benchmark object
type K8sAudit struct {
	Command     shell.Executor
	FailedTests []*models.AuditBench
	args        []string
}

//NewK8sAudit new audit object
func NewK8sAudit() *K8sAudit {
	return &K8sAudit{FailedTests: make([]*models.AuditBench, 0), Command: shell.NewShellExec()}
}

//Help return benchmark command help
func (bk K8sAudit) Help() string {
	return "-a , --audit run benchmark audit tests"
}

//Run execute benchmark command
func (bk *K8sAudit) Run(args []string) int {
	bk.args = args
	audit := models.Audit{}
	auditFiles, err := utils.GetK8sBenchAuditFiles()
	if err != nil {
		panic(fmt.Sprintf("failed to read audit files %s", err))
	}
	for _, auditFile := range auditFiles {
		err := json.Unmarshal([]byte(auditFile.Data), &audit)
		if err != nil {
			panic("Failed to unmarshal audit test json file")
		}
		for _, ac := range audit.Categories {
			bk.runTests(ac)
		}
	}
	reports.GenerateAuditReport(bk.FailedTests)

	return 0
}

func (bk *K8sAudit) runTests(ac models.Category) {
	for _, at := range ac.SubCategory.AuditTests {
		resArr := make([]string, 0)
		for index := range at.AuditCommand {
			res := bk.execCommand(at, index, resArr, make([]IndexValue, 0))
			resArr = append(resArr, res)
		}
		data := NewValidExprData(resArr, at)
		bk.evalExpression(data, make([]string, 0))
		if len(bk.args) == 1 && bk.args[0] != "report" {
			bk.printTestResults(data.atb)
		} else {
			bk.AddFailedMessages(data)
		}
	}
}

func (bk *K8sAudit) addDummyCommandResponse(at *models.AuditBench, index int) string {
	spExpr := utils.SeparateExpr(at.EvalExpr)
	for _, expr := range spExpr {
		if expr.Type == common.SingleValue {
			if !strings.Contains(expr.Expr, fmt.Sprintf("'$%d'", index)) {
				if strings.Contains(expr.Expr, fmt.Sprintf("$%d", index)) {
					return common.NotValidNumber
				}
			}
		}
	}
	return common.NotValidString
}

//AddFailedMessages add failed audit test to report data
func (bk *K8sAudit) AddFailedMessages(data ValidateExprData) {
	if data.atb.TestResult.NumOfSuccess != data.atb.TestResult.NumOfExec {
		bk.FailedTests = append(bk.FailedTests, data.atb)
	}
}

//IndexValue hold command index and result
type IndexValue struct {
	index int
	value string
}

func (bk *K8sAudit) execCommand(at *models.AuditBench, index int, prevResult []string, newRes []IndexValue) string {
	cmd := at.AuditCommand[index]
	paramArr, ok := at.CommandParams[index]
	if ok {
		for _, param := range paramArr {
			parmNum, err := strconv.Atoi(param)
			if err != nil {
				log.Console(fmt.Sprintf("failed to convert param for command %s", cmd))
				continue
			}
			if parmNum < len(prevResult) {
				n := prevResult[parmNum]
				if n == "[^\"]\\S*'\n" || n == "" || n == common.NotValidString {
					n = bk.addDummyCommandResponse(at, index)
				}
				newRes = append(newRes, IndexValue{index: parmNum, value: n})
			}
		}
		commandRes := make([]string, 0)
		bk.execCommandWithParams(newRes, len(newRes), make([]IndexValue, 0), len(newRes), cmd, &commandRes)
		sb := strings.Builder{}
		for _, cr := range commandRes {
			sb.WriteString(fmt.Sprintf("%s\n", cr))
		}
		return sb.String()
	}
	result, _ := bk.Command.Exec(cmd)
	if result.Stderr != "" {
		log.Console(fmt.Sprintf("Failed to execute command %s", result.Stderr))
	}
	return result.Stdout

}

func (bk *K8sAudit) execCommandWithParams(arr []IndexValue, index int, prevResHolder []IndexValue, origSize int, val string, resArr *[]string) {
	if len(arr) == 0 {
		return
	}
	sArr := strings.Split(arr[0].value, "\n")
	for _, a := range sArr {
		prevResHolder = append(prevResHolder, IndexValue{index: arr[0].index, value: a})
		bk.execCommandWithParams(arr[1:index], index-1, prevResHolder, origSize, val, resArr)
		if len(prevResHolder) == origSize {
			for _, param := range prevResHolder {
				if param.value == common.NotValidString || param.value == common.NotValidNumber || param.value == "" {
					*resArr = append(*resArr, param.value)
					break
				}
				cmd := strings.ReplaceAll(val, fmt.Sprintf("#%d", param.index), param.value)
				result, _ := bk.Command.Exec(cmd)
				if result.Stderr != "" {
					*resArr = append(*resArr, "")
					log.Console(fmt.Sprintf("Failed to execute command %s", result.Stderr))
				}
				*resArr = append(*resArr, result.Stdout)
			}
		}
		prevResHolder = prevResHolder[:len(prevResHolder)-1]
	}
}

func (bk *K8sAudit) printTestResults(at *models.AuditBench) {
	if at.TestResult.NumOfSuccess == at.TestResult.NumOfExec {
		log.Console(emoji.Sprintf(":check_mark_button: %s\n", at.Name))
	} else {
		log.Console(emoji.Sprintf(":cross_mark: %s\n", at.Name))
	}
}

func (bk *K8sAudit) evalExpression(ved ValidateExprData, combArr []string) {
	if len(ved.resultArr) == 0 {
		return
	}
	outputs := strings.Split(ved.resultArr[0], "\n")
	for _, o := range outputs {
		if len(o) == 0 && len(outputs) > 1 {
			continue
		}
		combArr = append(combArr, o)
		bk.evalExpression(ved.NextValidExprData(), combArr)
		if ved.origSize == len(combArr) {
			expr := ved.atb.Sanitize(combArr, ved.atb.EvalExpr)
			ved.atb.TestResult.NumOfExec++
			count, err := bk.evalCommandExpr(ved.atb, expr)
			if err != nil {
				log.Console(err.Error())
			}
			ved.atb.TestResult.NumOfSuccess += count
		}
		combArr = combArr[:len(combArr)-1]
	}

}

func (bk *K8sAudit) evalCommandExpr(at *models.AuditBench, expr string) (int, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, fmt.Errorf("failed to build evaluation command expr for\n %s", at.Name)
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, fmt.Errorf("failed to evaluate command expr for audit test %s", at.Name)
	}
	b, ok := result.(bool)
	if ok && b {
		return 1, nil
	}
	return 0, nil
}

//Synopsis for help
func (bk *K8sAudit) Synopsis() string {
	return bk.Help()
}
