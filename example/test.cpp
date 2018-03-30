#include "test.h"

testClass::testClass(std::string a):                                          
    _a(a)                                                          
{                                                                  
}                                                                  
                                                                            
std::string testClass::print()
{                                               
    return _a;                                                     
}
