/*
Copyright © 2021 Thomas Meitz

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// https://stackoverflow.com/questions/66067549/how-to-write-a-custom-error-reporter-in-go-target-of-antlr

package main

import (
	"fmt"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/oscarzhou/ksqldb-go/parser"
)

func main() {
	//----------| this is the error
	k := `select timestamptostring(windowstart,'yyyy-MM-dd HH:mm:ss','Europe/London') as windowStart, 
				  timestamptostring(windowend,'HH:mm:ss','Europe/London') as windowEnd, 
					dogSize, dogsCt 
	from1 dogs_by_size 
	where dog_size='large';`

	input := antlr.NewInputStream(k)
	upper := parser.NewUpperCaseStream(input)
	lexerErrors := &parser.KSqlErrorListener{}
	lexer := parser.NewKSqlLexer(upper)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexerErrors)

	stream := antlr.NewCommonTokenStream(lexer, 0)
	parserErrors := &parser.KSqlErrorListener{}
	p := parser.NewKSqlParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(parserErrors)

	antlr.ParseTreeWalkerDefault.Walk(&parser.BaseKSqlListener{}, p.Statements())
	fmt.Println(fmt.Sprintf("lexer has errors: %v", lexerErrors.HasErrors()))
	fmt.Println(fmt.Sprintf("parser error count: %v", lexerErrors.ErrorCount()))
	fmt.Println(fmt.Sprintf("parser has errors: %v", parserErrors.HasErrors()))
	fmt.Println(fmt.Sprintf("parser error count: %v", parserErrors.ErrorCount()))
	fmt.Println(parserErrors.Errors)
}
