Updates for v1.3.0 and v1.2.0:
 Velvet v1.3.0
  * Errors now messages to go with them which can be accessed using the `pusherr` instruction
  * All jump instructions now have a branch variant(e.g. `jt` -> `brt`), which allows the `ret` instruction to work
  * The `readNumber` function has been renamed to `readn` and three IO functions (`readt`, `readb`, and `readc`), have been added to go along with it; file reading/writing will be added in the next version

 Velvc v1.2.0
  * Fixed Generator struct not adding instructions and static data to the final executable (I freaking forgot to make a method a pointer receiver)
  * Added all new instruction variants
  * Added the `-v` flag to show the version

 Misc
  * Updated ImHex pattern script to include the new flags and variants
  * Fixed documentation mistake
