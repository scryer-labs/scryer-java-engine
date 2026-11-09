#include <dlfcn.h>
#include <stddef.h>

const char *java_engine_library_path(void) {
    Dl_info info;
    if (!dladdr((const void *)java_engine_library_path, &info)) {
        return NULL;
    }
    return info.dli_fname;
}
