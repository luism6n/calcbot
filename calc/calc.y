%{
package calc

import (
    "math"
)
%}

%union{
    val float64
    name string
}

%type <val> expr

%token <val> NUMBER
%token <val> IDENTIFIER
%token LOG
%token LOG10
%token LOG2
%token LN
%token POW
%token EXP

%left '+' '-'
%left '*' '/'

%%

prog : expr { yylex.(*calcLexer).result = $1 }
     | prog ';' expr { yylex.(*calcLexer).result = $3 }

expr : NUMBER { $$ = $1 }
     | expr '+' expr { $$ = $1 + $3 }
     | expr '-' expr { $$ = $1 - $3 }
     | expr '*' expr { $$ = $1 * $3 }
     | expr '/' expr { $$ = $1 / $3 }
     | '(' expr ')' { $$ = $2 }
     | LOG '(' expr ',' expr ')' { $$ = log($3, $5) }
     | LOG10 '(' expr ')' { $$ = log(10, $3) }
     | LOG2 '(' expr ')' { $$ = log(2, $3) }
     | LN '(' expr ')' { $$ = log(math.E, $3) }
     | POW '(' expr ',' expr ')' { $$ = pow($3, $5) }
     | EXP '(' expr ')' { $$ = exp($3) }
     | IDENTIFIER
     ;

%%
