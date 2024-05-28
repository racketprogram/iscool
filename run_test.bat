@echo off
setlocal enabledelayedexpansion

REM Initialize counters
set count=1
set fail_count=0

REM Loop to run tests 100 times
:loop
echo Running test iteration %count%
go test ./... -count=1 > test_output.txt 2>&1
IF %ERRORLEVEL% NEQ 0 (
    echo Test iteration %count% failed >> failed_tests.log
    type test_output.txt >> failed_tests.log
    set /a fail_count+=1
)
set /a count+=1
IF %count% LEQ 100 (
    goto loop
)

REM Summary
echo Test completed with %fail_count% failures out of 100 runs.
if %fail_count% NEQ 0 (
    echo Check failed_tests.log for more details.
)

endlocal
