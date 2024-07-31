echo off
FOR /F "eol=# tokens=*" %%i IN (.env) DO SET %%i
echo .env loaded
