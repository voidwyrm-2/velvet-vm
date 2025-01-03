#pragma author Nuclear Pasta
#pragma description Velvet Virtual Machine Bytecode Executable (VEX/cve/cvelv)

import std.mem;
import std.core;
import std.io;

std::core::set_endian(std::mem::Endian::Big); // Velvet bytecode uses big endian for all of it's numbers

struct InfoSection { // contains information about the file and program
    char signature[17] [[color("803D2F")]]; // signature 
    u8 vars [[color("00FF00")]]; // the amount of variables the program uses
    u16 dataAddress [[color("0000FF")]]; // the byte address of the Data Section
};

InfoSection info @ 0;

enum Opcode : u16 {
    NOP = 0,
    RET = 1,
    HALT = 2,
    CALL = 3,
    PUSH = 4,
    POP = 5,
    DUP = 6,
    SWAP = 7,
    ROT = 8,
    VAROP = 9,
    JUMP = 10
};

enum PushFlag : u8 {
    NUMBER = 0,
    BOOL = 1,
    STRING = 2,
    LIST = 3
};

enum SetFlag : u8 {
    SET = 0,
    GET = 1
};

enum JumpFlag : u8 {
    UNCOND_JUMP = 0,
    JUMP_IF_TRUE = 1,
    JUMP_IF_FALSE = 2,
    JUMP_IF_ERROR = 3,
    JUMP_IF_NOT_ERROR = 4
};

struct Instruction {
    Opcode opcode [[color("FF0000")]];
    if (this.opcode == Opcode::PUSH) {
        PushFlag flag [[color("D85656")]];
    } else if (this.opcode == Opcode::VAROP) {
        SetFlag flag [[color("D85656")]];
    } else if (this.opcode == Opcode::JUMP) {
        JumpFlag flag [[color("D85656")]];
    } else {
        u8 flag [[color("D85656")]];
    }
    
    if (this.opcode == Opcode::CALL || (this.opcode == Opcode::PUSH && (this.flag == PushFlag::STRING || this.flag == PushFlag::LIST))) {
        u16 argOne [[color("FF00FF"), name("address")]];
        u16 argTwo [[color("FF00FF"), name("length")]];
    } else {
        u16 argOne [[color("FF00FF")]];
        u16 argTwo [[color("FFFF00")]];
    }
};

struct VEX { // VevletVM Executable
    //InfoSection info;
    Instruction instructions[while($ < info.dataAddress)];
    u8 data[while(!std::mem::eof())] [[color("5FCDE4")]]; // contains various static data
};

VEX vex @ 20;