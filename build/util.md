#Generate zip without __MACOSX files
zip -r dir.zip . -x '**/.*' -x '**/__MACOSX'