# Build Your Own Calculator
This challenge is to build your own calculator.

## The challenge
To complete this challenge you’ll need to parse a mathematical expression and then perform the relevant calculations before returning the answer to the user.

For example, the user will be able to input: 2 * 3 + 4 and get back 10, or input 10 / (6 - 1) and get back 2.

Completing this challenge will give you the chance to make use of the stack data structure in a real-world application.

## Run locally 
```bash
go run . <expression>
# eg
go run . '1 + 2'
```

## Solution
The solution uses reverse polish notation to solve the problem. At first tokenize the given input. Then create the reverse polish notation of the tokens using shunting yard algorithm. Use the RPN to solve the equation

### Operators
Create a operator map to identify the operators. The operators have precedence that are used during the shunting yard algorithm.

### Tokenization
Read through the expression and parse it based on the current character. Parse the expression carefully to identify the numbers

### Shunting yard algorithm
Read more about the algorithm from [here](https://www.andreinc.net/2010/10/05/converting-infix-to-rpn-shunting-yard-algorithm). 

### Evaluation of the RPN 
Read through the tokens generated using the shunting yard algorithm and evaluate them using a stack. If the current token is an operator, pop two element from stack, apply the operation on them and push the result in the stack.

### Identifying malformed expression
In every stage of the evaluation, perform checks for malformed expression and output appropriate exception messages.