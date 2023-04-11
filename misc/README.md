![Our logo] (logo.svg)

$env:ANDROID_SDK_ROOT="C:\Users\t-ethanperry\AppData\Local\Android\Sdk"

Notes:
1. If our background tasks are expected to take a REALLY long time, we may consider using https://docs.celeryproject.org/en/stable/getting-started/introduction.html (Celery) or other distributed task queue. These are useful