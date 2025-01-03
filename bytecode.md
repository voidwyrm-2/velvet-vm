# Velvet Bytecode

Velvet bytecode consists of three sections: Info, Instruction, and Data

## Info

The Info section consists of 32 bytes

- The first 17 spell out "Velvet Scarlatina"
- The 18th byte tells the VM how many variables are used during the program
- The last three bytes tells the VM the position of where the Data section starts

## Instruction

The Instruction section contains all of the instructions that the VM runs

Each instruction is made up of 7 bytes, split into three to four sections, like so
```
00 0 00 00
00 0 0000
```
The first section of two bytes is the opcode

The second section of one byte is the instruction flag, it varies between instructions

The third and fourth sections of two bytes are either read separately or together as one section, depending on the instruction

## Data

The Data section contains certain static data values and function names used in the program, arrayed back-to-back

How the Data section is used is that instructions that use this section (`push` and `call` at the time of writing)
have an address and a data length in place of arguments, which are then used to get the value from the data section

### Data Formating

* Strings are just sequences of bytes; they are **not** null-terminated
* Lists are sequences of items composed of 5 bytes each; one for type, two for address, two for length; lists can strings, numbers, booleans, or other lists
