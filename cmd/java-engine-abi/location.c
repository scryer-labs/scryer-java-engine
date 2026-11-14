#ifdef _WIN32
#include <windows.h>

static _Thread_local wchar_t library_path[32768];
static _Thread_local char utf8_path[131072];

const char *java_engine_library_path(void) {
    HMODULE module;
    if (!GetModuleHandleExW(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS | GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT,
                            (LPCWSTR)java_engine_library_path, &module)) return NULL;
    DWORD length = GetModuleFileNameW(module, library_path, 32768);
    if (length == 0 || length >= 32768) return NULL;
    if (!WideCharToMultiByte(CP_UTF8, 0, library_path, -1, utf8_path, sizeof(utf8_path), NULL, NULL)) return NULL;
    return utf8_path;
}
#else
#include <dlfcn.h>
#include <stddef.h>

const char *java_engine_library_path(void) {
    Dl_info info;
    if (!dladdr((const void *)java_engine_library_path, &info)) {
        return NULL;
    }
    return info.dli_fname;
}
#endif
