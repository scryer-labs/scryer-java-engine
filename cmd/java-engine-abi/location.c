#include <stdlib.h>

#ifdef _WIN32
#include <windows.h>

char *java_engine_library_path(void) {
    HMODULE module;
    if (!GetModuleHandleExW(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS | GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT,
                            (LPCWSTR)java_engine_library_path, &module)) return NULL;

    wchar_t path[32768];
    DWORD length = GetModuleFileNameW(module, path, sizeof(path) / sizeof(path[0]));
    if (length == 0 || length >= sizeof(path) / sizeof(path[0])) return NULL;

    int size = WideCharToMultiByte(CP_UTF8, 0, path, -1, NULL, 0, NULL, NULL);
    if (size == 0) return NULL;
    char *utf8 = malloc(size);
    if (utf8 == NULL) return NULL;
    if (!WideCharToMultiByte(CP_UTF8, 0, path, -1, utf8, size, NULL, NULL)) {
        free(utf8);
        return NULL;
    }
    return utf8;
}
#else
#include <dlfcn.h>
#include <string.h>

char *java_engine_library_path(void) {
    Dl_info info;
    if (!dladdr((const void *)java_engine_library_path, &info)) {
        return NULL;
    }
    return strdup(info.dli_fname);
}
#endif
