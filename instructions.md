# Instructions

**Note:** "stack item" refers to the item at the top of the stack

0. nop
1. ret
2. halt (int8)
3. call (fn)
    **Flags**<br>
    default: Treats the instruction arguments as an address and length for the name of a function
    1. Ignores the instruction arguments and instead pops off a function from the stack and runs it
4. push (literal)

    **Flags**<br>
    default: Treats the instruction arguments as a number
    1. Treats the instruction arguments as a bool
    2. Treats the instruction arguments as an address and length for a string
    3. Treats the instruction arguments as an address and length for a list
    4. Treats the instruction arguments as an address and length for the name of a function
5. pop
6. dup
7. swap
8. rot
9. set (int16)

    **Flags**<br>
    default: Instruction sets the variable of the given index to the stack item
    1. Instruction gets the variable of the given index and pushes its value onto the stack
10. j (label)

    **Flags**<br>
    default: Instruction acts like an unconditional jump
    1. Instruction only jumps if the stack item is `true`
    2. Instruction only jumps if the stack item is `false`
    3. Instruction only jumps if the error flag is true (a previous function call errored out)
    4. Instruction only jumps if the error flag is false (a previous function call did not error out)

# Function Instructions

These instructions don't actually exist, but are converted into function calls during the compilation process

* `error` -> `errflag = true`
* `reset` -> `errflag = false`
* `eq` -> `y, x = pop(), pop(); push(x == y)`
* `neq` -> `y, x = pop(), pop(); push(x != y)`
* `not` -> `x = pop(); push(!x)`
* `lt` -> `y, x = pop(), pop(); push(x < y)`
* `gt` -> `y, x = pop(), pop(); push(x > y)`
* `lte` -> `y, x = pop(), pop(); push(x <= y)`
* `gte` -> `y, x = pop(), pop(); push(x >= y)`
* `add` -> `y, x = pop(), pop(); push(x + y)`
* `sub` -> `y, x = pop(), pop(); push(x - y)`
* `mul` -> `y, x = pop(), pop(); push(x * y)`
* `div` -> `y, x = pop(), pop(); push(x / y)`
* `pow` -> `y, x = pop(), pop(); push(pow(x, y))`
* `log` -> `y, x = pop(), pop(); push(log(x, y))`
* `neg` -> `x = pop(); push(x!)`
* `and` -> `y, x = pop(), pop(); push(x & y)`
* `or` -> `y, x = pop(), pop(); push(x | y)`
* `xor` -> `y, x = pop(), pop(); push(x ^ y)`
