%{
    package calc
%}

%union{
    val float64
}

%token NUMBER

%%

expr : NUMBER ;

%%
