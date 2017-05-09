@0xa49bdccd9f1db7ca;
using Go = import "go.capnp";
$Go.package("protocol");
$Go.import("testpkg");


struct ConcatReplyCapn { 
   retCode  @0:   Int8; 
   val      @1:   Text; 
} 

struct ConcatRequestCapn { 
   userId  @0:   Int64; 
   str1    @1:   Text; 
   str2    @2:   Text; 
} 

struct MatchReplyCapn { 
   retCode  @0:   Int8; 
} 

struct MatchRequestCapn { 
   gameGroup  @0:   Text; 
   gameDb     @1:   Text; 
   serverId   @2:   Int64; 
   userId     @3:   Int64; 
   score      @4:   Int64; 
} 

struct SumReplyCapn { 
   retCode  @0:   Int8; 
   val      @1:   Int64; 
} 

struct SumRequestCapn { 
   userId  @0:   Int64; 
   num1    @1:   Int64; 
   num2    @2:   Int64; 
} 

##compile with:

##
##
##   capnp compile -ogo ./schema.capnp

