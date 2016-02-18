# Retry

## About

Retry is a [Go](https://golang.org/) utility package to to facilitate operations
that need to be retried when they fail.

## Example

Retry an operation as many times as necessary until it completes successfully.

```golang
options := retry.Options{
    Do: func() (retry.Result, error) {
        if err := DoSomething(); err != nil {
            return retry.Continue, err
        } else {
            return retry.Stop, nil
        }
    },
}

if err := retry.Do(options); err != nil {
    log.Println("retry failed:", err)
}
```

## Features

*Limit the number of retry attempts (default forever)*

```golang
options := retry.Options {
    MaxAttempts: 10,
}
```

*Exponential backoff with max delay (default 30 seconds)*

```golang
options := retry.Options {
    MaxDelay: time.Second * 30
}
```

*Deadline based retry window (default forever)*

```golang
options := retry.Options {
    Deadline: time.Now().UTC().Add(time.Minute * 5)
}
```

*Optional cancellation (work in progress)*

```golang
cancelled := make(chan interface{})

options := retry.Options{
    Do: func() (retry.Result, error) {
        if err := DoSomething(cancelled); err != nil {
            return retry.Continue, err
        } else {
            return retry.Stop, nil
        }
    },
    Cancel: func() error {
        close(cancelled)
        return nil
    }
}

retry := retry.New(options)

go retry.Do()

if err := retry.Cancel(); err != nil {
    log.Println("retry cancel failed:", err)
}
```

## LICENSE

```
The MIT License (MIT)

Copyright (c) 2016 Mediafly, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
