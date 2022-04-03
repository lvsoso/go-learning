#include <iostream>

using namespace std;

extern "C" {
    void test()
    {
        cout << "this is a test" << endl;
    }
}