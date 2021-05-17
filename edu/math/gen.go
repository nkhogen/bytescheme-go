package math

import (
	"time"
	"math/rand"
	"bytes"
	//"fmt"
	"strconv"
	"strings"
	"io/ioutil"
)


type Config struct {
	MultConf *MultConfig
	DivConf *DivConfig
	SubConf *SubConfig
	Cols int
}

type MultConfig struct {
	FirstLen int
	SecondLen int
	Size int
}

type DivConfig struct {
	QuotientLen int
	DivisorLen int
	Size int
}

type SubConfig struct {
	SubtractorLen int
	ResultLen int
	Size int
}

type BinaryExpression struct {
	FirstOperand string
	SecondOperand string
	
}

var (
	random = rand.New(rand.NewSource(time.Now().Unix()))
)

func init() {

}

type ProblemGenerator struct {
	Conf *Config
}

func generateNum(start, digitLen int) string {
	buff := &bytes.Buffer{}
	for i := 0; i < digitLen; i++ {
		n := start + random.Intn(10-start)
		if buff.Len() > 0 {
			buff.WriteByte(' ')
		}
		buff.WriteByte(byte('0'+n))
	}
	return buff.String()
}

func (gen *ProblemGenerator) GenerateHTMLFile(filename string) error {
	html := gen.GenerateHTML()
	err := ioutil.WriteFile(filename, html, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (gen *ProblemGenerator) GenerateHTML() []byte {
	serial := 1
	buff := &bytes.Buffer{}
	buff.WriteString("<!DOCTYPE html>\n")
	buff.WriteString("<html>\n<head>\n")
	buff.WriteString("<style>\n")
	buff.WriteString(`table, th, td {
		border: 0px solid black;
		padding: 5px;
		text-align: right;
		vertical-align: top;
		padding-right: 60px;
		padding-bottom: 40px;
	}
	table {
		border-spacing: 15px;
	}
	</style>`)
	buff.WriteString("\n</head>\n<body>\n")
	buff.WriteString("<center>Name: Punya Naorem<br>Date: <script>document.write(new Date().toLocaleDateString());</script></center>")
	buff.WriteString("<table style=\"width:100%\">")
	gen.GenerateSub(buff, &serial)
	gen.GenerateMult(buff, &serial)
	gen.GenerateDiv(buff, &serial)
	buff.WriteString("</table>\n</head>\n</body>")
	return buff.Bytes()
}

func (gen *ProblemGenerator) GenerateMult(buff *bytes.Buffer, serial *int) {
	config := gen.Conf.MultConf
	for i := 0; i < config.Size; i++ {
		buff.WriteString("<tr>")
		for j := 0; j < gen.Conf.Cols; j++ {
			if j > 0 {
				buff.WriteString("<td></td>")
			}
			buff.WriteString("<td>")
			buff.WriteString(strconv.Itoa(*serial))
			buff.WriteString(".</td>")
			first := generateNum(1, config.FirstLen)
			second :=  generateNum(1, config.SecondLen)
			buff.WriteString("<td>")
			buff.WriteString(first)
			buff.WriteString("<br>")
			buff.WriteString("x")
			buff.WriteString("<br>")
			buff.WriteString(second)
			buff.WriteString("<br>")
			buff.WriteString("--------")
			for k := 0; k < config.SecondLen; k++ {
				buff.WriteString("<br>")
			}
			buff.WriteString("</td>")
			*serial = *serial + 1
		}
		buff.WriteString("</tr>\n")
	}

}

func (gen *ProblemGenerator) GenerateDiv(buff *bytes.Buffer, serial *int) {
	config := gen.Conf.DivConf
	for i := 0; i < config.Size; i++ {
		buff.WriteString("<tr>")
		for j := 0; j < gen.Conf.Cols; j++ {
			if j > 0 {
				buff.WriteString("<td></td>")
			}
			buff.WriteString("<td>")
			buff.WriteString(strconv.Itoa(*serial))
			buff.WriteString(".</td>")
			quotientStr := generateNum(1, config.QuotientLen)
			quotient, _ := strconv.Atoi(strings.ReplaceAll(quotientStr, " ", ""))
			divisorStr := generateNum(2, config.DivisorLen)
			divisor, _ := strconv.Atoi(strings.ReplaceAll(divisorStr, " ", ""))
			dividend := divisor * quotient
			buff.WriteString("<td>")
			buff.WriteString(strconv.Itoa(dividend))
			buff.WriteString("&nbsp;/&nbsp;")
			buff.WriteString(strconv.Itoa(divisor))
			buff.WriteString("&nbsp;=")
			buff.WriteString("</td>")
			*serial = *serial + 1
		}
		buff.WriteString("</tr>\n")
	}
}

func (gen *ProblemGenerator) GenerateSub(buff *bytes.Buffer, serial *int) {
	config := gen.Conf.SubConf
	for i := 0; i < config.Size; i++ {
		buff.WriteString("<tr>")
		for j := 0; j < gen.Conf.Cols; j++ {
			if j > 0 {
				buff.WriteString("<td></td>")
			}
			buff.WriteString("<td>")
			buff.WriteString(strconv.Itoa(*serial))
			buff.WriteString(".</td>")
			subtractorStr := generateNum(1, config.SubtractorLen)
			subtractor, _ := strconv.Atoi(strings.ReplaceAll(subtractorStr, " ", ""))
			resultStr :=  generateNum(1, config.ResultLen)
			result, _ := strconv.Atoi(strings.ReplaceAll(resultStr, " ", ""))
			firstStr := strconv.Itoa(subtractor + result)
			firstStr = strings.Join(strings.Split(firstStr, ""), " ")
			buff.WriteString("<td>")
			buff.WriteString(firstStr)
			buff.WriteString("<br>")
			buff.WriteString("-")
			buff.WriteString("<br>")
			buff.WriteString(subtractorStr)
			buff.WriteString("<br>")
			buff.WriteString("--------")
			buff.WriteString("<br>")
			buff.WriteString("</td>")
			*serial = *serial + 1
		}
		buff.WriteString("</tr>\n")
	}

}