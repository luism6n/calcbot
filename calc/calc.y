%{
    package calc
%}

%union{
    val int
}

%token NUMBER

%%

expr : NUMBER ;

%%
