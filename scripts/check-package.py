"""Exercise the packaged v1 ABI with a bundled Java runtime."""
import ctypes
import json
import os
import sys
import tempfile
from pathlib import Path

package = Path(sys.argv[1]).resolve()
manifest = json.loads((package / "engine.json").read_text())
library = ctypes.CDLL(str(package / manifest["library"]))


class View(ctypes.Structure):
    _fields_ = [("data", ctypes.c_char_p), ("size", ctypes.c_size_t)]


class Buffer(ctypes.Structure):
    _fields_ = [("data", ctypes.c_void_p), ("size", ctypes.c_size_t)]


class Result(ctypes.Structure):
    _fields_ = [("status", ctypes.c_int32), ("response", Buffer)]


library.execute_v1.argtypes = [View]
library.execute_v1.restype = Result
library.cancel_v1.argtypes = [View]
library.cancel_v1.restype = ctypes.c_bool
library.release_buffer_v1.argtypes = [Buffer]
library.release_buffer_v1.restype = None
with tempfile.TemporaryDirectory() as directory:
    root = Path(directory).resolve()
    os.environ["JAVA_HOME"] = str(root / "missing")
    os.environ.pop("SCRYER_ANALYZER_JAVA_HOME", None)
    os.environ.pop("SCRYER_JAVA_CLASSPATH", None)
    (root / "Main.java").write_text("class Main {}\n")
    (root / "pom.xml").write_text("<project><modelVersion>4.0.0</modelVersion><groupId>test</groupId><artifactId>check</artifactId><version>1</version></project>")
    request = json.dumps({
        "protocolVersion": 1, "requestId": "package-check", "stack": "java",
        "operation": "scan", "scan": {"repository": str(root), "facts": ["dependencies"], "static": True},
    }).encode()
    result = library.execute_v1(View(request, len(request)))
    try:
        response = json.loads(ctypes.string_at(result.response.data, result.response.size))
        assert result.status == 0 and response["error"] is None, response
        assert "dependencies" in response["report"]["facts"]
    finally:
        library.release_buffer_v1(result.response)
