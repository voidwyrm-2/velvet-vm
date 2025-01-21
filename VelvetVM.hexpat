#pragma author Nuclear Pasta
#pragma description Velvet Virtual Machine Bytecode Executable (VEX/cve/cvelv)

import std.mem;
import std.core;
import std.io;

std::core::set_endian(std::mem::Endian::Big); // Velvet bytecode uses big endian for all of its numbers

struct InfoSection { // contains information about the file and program
    char signature[17] [[color("803D2F")]]; // signature
    u8 programFlags [[color("AAFF00")]];
    u16 vars [[color("00FF00")]]; // the amount of variables the program uses
    u32 dataAddress [[color("0000FF")]]; // the byte address of the Data Section
    u32 programEntry [[color("00AAFF")]];
    u32 reserved [[color("000000")]];
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
    JUMPBRANCH = 10
};

enum PushFlag : u8 {
    NUMBER = 0,
    BOOL = 1,
    STRING = 2,
    LIST = 3,
    FUNCTION = 4,
    ERRREG = 5
};

enum CallFlag : u8 {
    FUNCTION_IN_LIBRARY = 0,
    FUNCTION_ON_STACK = 1
};

enum SetFlag : u8 {
    SET = 0,
    GET = 1
};

enum JumpBranchFlag : u8 {
    UNCOND_JUMP = 0x0000,
    UNCOND_BRANCH = 0x1000,
    JUMP_IF_TRUE = 0x0001,
    BRANCH_IF_TRUE = 0x1001,
    JUMP_IF_FALSE = 0x0010,
    BRANCH_IF_FALSE = 0x1010,
    JUMP_IF_ERROR = 0x0011,
    BRANCH_IF_ERROR = 0x1011,
    JUMP_IF_NOT_ERROR = 0x0100,
    BRANCH_IF_NOT_ERROR = 0x1100
};


struct Instruction {
    Opcode opcode [[color("FF0000")]];
    if (this.opcode == Opcode::PUSH) {
        PushFlag flag [[color("D85656")]];
    } else if (this.opcode == Opcode::CALL) {
        CallFlag flag [[color("D85656")]];
    } else if (this.opcode == Opcode::VAROP) {
        SetFlag flag [[color("D85656")]];
    } else if (this.opcode == Opcode::JUMPBRANCH) {
        JumpBranchFlag flag [[color("D85656")]];
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

VEX vex @ 32;
