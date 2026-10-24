#!/usr/bin/env python3
import hashlib
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
SWAY_LOGOS = {
    "assets/branding/sway_s.png": "3e8802a71c97374dffccf76f6df25348aff3b47d2dafa79685ef1f0c5ec4a7c6",
    ".github/assets/sway-logo-black.svg": "9aea2f5c3d7a3490b31e9340419e17083ce2f1ed23d8edf0088533d77e837851",
    ".github/assets/sway-logo-white.svg": "3963583ee0b5f077985b9a872f372aa3efcfdb81dc427e7269ca3b1edd7c3bf4",
}
SWAY_LOGO_REFERENCES = [
    ("src/fr/sway/accueil/TopBar.java", '"branding/sway_s.png"'),
    ("README.md", 'srcset=".github/assets/sway-logo-white.svg"'),
    ("README.md", 'src=".github/assets/sway-logo-black.svg"'),
]
MAX_JAVA_FILE_LINES = 200
MAX_JAVA_METHOD_LINES = 50
MAX_UTILITY_METHOD_LINES = 25
MAX_UTILITY_METHODS = 5
JAVA_TOP_LEVEL_TYPE = re.compile(r"^(public |final |abstract )*(class|interface|enum) ")
JAVA_WILDCARD_IMPORT = re.compile(r"^import (static )?[\w.]+\.\*;")
JAVA_CONDITION_ASSIGNMENT = re.compile(r"\b(if|while)\s*\(.*\(\s*\w+\s*=[^=]")
DASH_GLYPHS = "[\\u00b7\\u2013\\u2014]"
FORBIDDEN_DASHES = re.compile(DASH_GLYPHS + "|\\S - \\S")
FORBIDDEN_DOC_DASHES = re.compile(DASH_GLYPHS)
JAVA_STRING = re.compile(r'"(?:[^"\\]|\\.)*"')
JAVA_METHOD_START = re.compile(r"^    \S.*\)\s*(throws [\w., ]+)?\{$")
PACKAGE_LINE = re.compile(r"^[a-z][\w.]*\*? \| \S.*$")
MAX_GO_FILE_LINES = 200
MAX_GO_TEST_FILE_LINES = 400
MAX_GO_FUNCTION_LINES = 100
MAX_GO_LIBRARY_FUNCTIONS = 5
GO_STRING = re.compile(r'"(?:[^"\\]|\\.)*"|`[^`]*`')
GO_ERROR_IN_CONDITION = re.compile(r"^\s*if .*:=.*;.*\berr\b")
GO_EXEC_OWNER = "bridge/internal/tv/adb.go"
TEXT_GLOBS = ["src/**/*.java", "scripts/*.sh", "scripts/*.py", "build.sh", "Makefile",
              "config/**/*.txt", "*.txt", "*.md", ".github/**/*.md", ".github/**/*.yml", "AndroidManifest.xml",
              "bridge/**/*.go", "bridge/go.mod", "bridge/deploy/*", "docs/*.md"]

errors = []

def report(path, line, code, message):
    errors.append(f"{path.relative_to(ROOT)}:{line}: [{code}] {message}")

def check_common(path, text, fix):
    lines = text.split("\n")
    if fix:
        fixed = re.sub(r"\n{3,}", "\n\n", "\n".join(l.rstrip() for l in lines)).rstrip("\n") + "\n"
        path.write_text(fixed)
        text, lines = fixed, fixed.split("\n")
    for number, line in enumerate(lines, 1):
        if line != line.rstrip():
            report(path, number, "N-ALL-01", "trailing whitespace")
        if number > 1 and line == "" and lines[number - 2] == "" and number < len(lines):
            report(path, number, "N-ALL-02", "consecutive empty lines")
    if not text.endswith("\n") or text.endswith("\n\n"):
        report(path, len(lines), "N-ALL-03", "file must end with exactly one newline")
    return lines

def check_user_text(path, number, text):
    if FORBIDDEN_DASHES.search(text):
        report(path, number, "N-ALL-04", "middle dot, en/em dash or ' - ' separator in user-facing text")

def is_utility_class(path, lines):
    return any(line.strip() == f"private {path.stem}() {{" for line in lines)

def check_documentation(path, lines):
    for number, line in enumerate(lines, 1):
        if FORBIDDEN_DOC_DASHES.search(line):
            report(path, number, "N-ALL-04", "middle dot or en/em dash in documentation")

def check_java_lines(path, lines):
    for number, line in enumerate(lines, 1):
        code = JAVA_STRING.sub('""', line)
        if "//" in code or "/*" in code:
            report(path, number, "N-JAVA-03", "comment in Java code")
        if JAVA_WILDCARD_IMPORT.match(code):
            report(path, number, "N-JAVA-05", "wildcard import")
        if JAVA_CONDITION_ASSIGNMENT.search(code):
            report(path, number, "N-JAVA-06", "assignment inside a condition")
        for literal in JAVA_STRING.findall(line):
            check_user_text(path, number, literal)

def check_java_methods(path, lines, utility):
    limit = MAX_UTILITY_METHOD_LINES if utility else MAX_JAVA_METHOD_LINES
    method_start, method_count = None, 0
    for number, line in enumerate(lines, 1):
        if JAVA_METHOD_START.match(line) and method_start is None:
            method_start = number
        elif line == "    }" and method_start is not None:
            method_count += 1
            if number - method_start + 1 > limit:
                report(path, method_start, "N-JAVA-02", f"method longer than {limit} lines")
            method_start = None
    if utility and method_count - 1 > MAX_UTILITY_METHODS:
        report(path, 1, "N-JAVA-07", f"more than {MAX_UTILITY_METHODS} methods in a utility class")

def check_java(path, lines):
    if len(lines) - 1 > MAX_JAVA_FILE_LINES:
        report(path, len(lines), "N-JAVA-01", f"file longer than {MAX_JAVA_FILE_LINES} lines")
    if sum(1 for line in lines if JAVA_TOP_LEVEL_TYPE.match(line)) > 1:
        report(path, 1, "N-JAVA-04", "more than one top-level class")
    check_java_lines(path, lines)
    check_java_methods(path, lines, is_utility_class(path, lines))

def check_go_lines(path, lines):
    for number, line in enumerate(lines, 1):
        code = GO_STRING.sub('""', line)
        if ("//" in code and not code.lstrip().startswith("//go:")) or "/*" in code:
            report(path, number, "N-GO-07", "comment in Go code")
        if GO_ERROR_IN_CONDITION.match(code):
            report(path, number, "N-GO-03", "assign and test an error on two lines")
        if "exec.Command" in code and path.relative_to(ROOT).as_posix() != GO_EXEC_OWNER:
            report(path, number, "N-SEC-02", f"processes may only be started from {GO_EXEC_OWNER}")

def check_go_functions(path, lines):
    start, count = None, 0
    for number, line in enumerate(lines, 1):
        if line.startswith("func ") and line.endswith("{"):
            start, count = number, count + 1
        elif line == "}" and start is not None:
            if number - start + 1 > MAX_GO_FUNCTION_LINES:
                report(path, start, "N-GO-02", f"function longer than {MAX_GO_FUNCTION_LINES} lines")
            start = None
    library = "/internal/" in path.as_posix() and not path.name.endswith("_test.go")
    if library and count > MAX_GO_LIBRARY_FUNCTIONS:
        report(path, 1, "N-GO-06", f"more than {MAX_GO_LIBRARY_FUNCTIONS} functions in a library file")

def check_go(path, lines):
    limit = MAX_GO_TEST_FILE_LINES if path.name.endswith("_test.go") else MAX_GO_FILE_LINES
    if len(lines) - 1 > limit:
        report(path, len(lines), "N-GO-01", f"file longer than {limit} lines")
    if path.name.startswith("tool_") and not path.name.endswith("_test.go"):
        if not path.with_name(path.stem + "_test.go").exists():
            report(path, 1, "N-GO-05", "every MCP tool needs its own test file")
    check_go_lines(path, lines)
    check_go_functions(path, lines)

def check_security(path, lines):
    if path.name == "go.mod" and any(line.startswith("require") for line in lines):
        report(path, 1, "N-SEC-03", "the bridge must not depend on third-party modules")
    if path.name == "AndroidManifest.xml" and any("<uses-permission" in line for line in lines):
        report(path, 1, "N-SEC-01", "the launcher must not request any Android permission")

def check_shell(path, lines, is_makefile):
    if not is_makefile and lines[:2] != ["#!/bin/bash", "set -euo pipefail"]:
        report(path, 1, "N-SH-01", "script must start with #!/bin/bash and set -euo pipefail")
    for number, line in enumerate(lines, 1):
        if line.lstrip().startswith("#") and not (number == 1 and line.startswith("#!")):
            report(path, number, "N-SH-02", "comment in shell or Makefile")

def check_package_list(path, lines):
    for number, line in enumerate(lines, 1):
        if line and not PACKAGE_LINE.match(line):
            report(path, number, "N-PKG-01", "expected 'package.name | Description'")
        if " | " in line:
            check_user_text(path, number, line.split(" | ", 1)[1])

def check_sway_logos():
    for relative, expected in SWAY_LOGOS.items():
        logo = ROOT / relative
        if not logo.exists():
            report(logo.parent, 0, "N-SWAY-01", f"Sway logo {logo.name} was removed")
        elif hashlib.sha256(logo.read_bytes()).hexdigest() != expected:
            report(logo, 0, "N-SWAY-01", "Sway logo was modified")
    for source, reference in SWAY_LOGO_REFERENCES:
        if reference not in (ROOT / source).read_text():
            report(ROOT / source, 0, "N-SWAY-02", "Sway logo is no longer displayed")

def files(targets):
    roots = [ROOT / t for t in targets] or [ROOT]
    for glob in TEXT_GLOBS:
        for path in sorted(ROOT.glob(glob)):
            if path.is_file() and any(path == r or r in path.parents for r in roots):
                yield path

def main():
    fix = "--fix" in sys.argv
    targets = [a for a in sys.argv[1:] if not a.startswith("--")]
    for path in files(targets):
        lines = check_common(path, path.read_text(), fix)
        if path.suffix == ".java":
            check_java(path, lines)
        elif path.suffix == ".go":
            check_go(path, lines)
        elif path.name in ("go.mod", "AndroidManifest.xml"):
            check_security(path, lines)
        elif path.suffix == ".sh" or path.name == "Makefile":
            check_shell(path, lines, path.name == "Makefile")
        elif path.parent.name == "packages" and path.suffix == ".txt":
            check_package_list(path, lines)
        elif path.suffix in (".md", ".txt"):
            check_documentation(path, lines)
    check_sway_logos()
    print("\n".join(errors) if errors else "Norminette: OK")
    return 1 if errors else 0

if __name__ == "__main__":
    sys.exit(main())
