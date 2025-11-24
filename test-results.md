k6-test.js:
тесты используют заранее созданные данные: 20 команди 200 пользователей (создаются data-for-loadtest.js)

Тест №1

операцими:
3 виртуальных пользователя с задержкой 0.2 секунды выполняют случайный запрос:

30% - получение команды 
30% - получение ревью
20% - изменение активности пользователя
20% - мерж PR

результаты:
RPS: 7.0
latency p95: 6.38ms   
ошибки: 0.00% 
проверки: 100% успешно 

вывод k6:
  █ THRESHOLDS

    http_req_duration
    ✓ 'p(95)<300' p(95)=6.38ms

    http_req_failed
    ✓ 'rate<0.001' rate=0.00%

    http_reqs
    ✓ 'rate>5' rate=7.048105/s


  █ TOTAL RESULTS

    checks_total.......: 439     7.048105/s
    checks_succeeded...: 100.00% 439 out of 439
    checks_failed......: 0.00%   0 out of 439

    ✓ get team
    ✓ get reviews
    ✓ merge PR
    ✓ set active

    HTTP
    http_req_duration..............: avg=4.47ms   min=2.35ms   med=4.52ms   max=8.64ms   p(90)=6.1ms    p(95)=6.38ms
      { expected_response:true }...: avg=4.47ms   min=2.35ms   med=4.52ms   max=8.64ms   p(90)=6.1ms    p(95)=6.38ms
    http_req_failed................: 0.00%  0 out of 439
    http_reqs......................: 439    7.048105/s

    EXECUTION
    iteration_duration.............: avg=205.36ms min=202.72ms med=205.47ms max=210.24ms p(90)=207.08ms p(95)=207.38ms
    iterations.....................: 439    7.048105/s
    vus............................: 2      min=1        max=2
    vus_max........................: 3      min=3        max=3

    NETWORK
    data_received..................: 1.5 MB 24 kB/s
    data_sent......................: 61 kB  985 B/s

running (1m02.3s), 0/3 VUs, 439 complete and 0 interrupted iterations
default ✓ [ 100% ] 0/3 VUs  1m0s

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Тест №2

нагрузка:
30 → 60 → 90 виртуальных пользователей
задержка 0.05 секунды

операции:
40% - получение команды
30% - получение ревью 
30% - изменение активности пользователя

результаты:
RPS: 490
ошибки: 0.00%  
проверки: 100% успешно 
latency p95: 143ms

вывод k6:

  █ THRESHOLDS

    http_req_duration
    ✓ 'p(95)<300' p(95)=143.03ms

    http_req_failed
    ✓ 'rate<0.01' rate=0.00%

    http_reqs
    ✓ 'rate>200' rate=490.494535/s


  █ TOTAL RESULTS

    checks_total.......: 35446   490.494535/s
    checks_succeeded...: 100.00% 35446 out of 35446
    checks_failed......: 0.00%   0 out of 35446

    ✓ get team
    ✓ get reviews
    ✓ set active

    HTTP
    http_req_duration..............: avg=50.98ms min=1.99ms  med=35.93ms max=743.02ms p(90)=111.62ms p(95)=143.03ms
      { expected_response:true }...: avg=50.98ms min=1.99ms  med=35.93ms max=743.02ms p(90)=111.62ms p(95)=143.03ms
    http_req_failed................: 0.00%  0 out of 35446
    http_reqs......................: 35446  490.494535/s

    EXECUTION
    iteration_duration.............: avg=101.1ms min=52.51ms med=87.33ms max=487.33ms p(90)=162.31ms p(95)=192.78ms
    iterations.....................: 35446  490.494535/s
    vus............................: 31     min=3          max=89
    vus_max........................: 90     min=90         max=90

    NETWORK
    data_received..................: 259 MB 3.6 MB/s
    data_sent......................: 4.7 MB 65 kB/s

running (1m12.3s), 00/90 VUs, 35446 complete and 0 interrupted iterations
default ✓ [ 100% ] 00/90 VUs  1m10s

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Тест №3

нагрузка:
100 → 300 → 500 виртуальных пользователей
задержка 0.005 секунды 

операции:
25% - получение команды
25% - получение ревью 
25% - создание новых PR
25% - изменение активности + мерж PR

результаты:
RPS: 506
ошибки: 0.22%
проверки: 99.78% успешно 
latency p95: 1.84s 

вывод k6:

  █ THRESHOLDS

    http_req_duration
    ✓ 'p(95)<2000' p(95)=1.84s

    http_req_failed
    ✓ 'rate<0.5' rate=0.22%

    http_reqs
    ✗ 'rate>1000' rate=506.682287/s


  █ TOTAL RESULTS

    checks_total.......: 17999  405.341326/s
    checks_succeeded...: 99.78% 17961 out of 17999
    checks_failed......: 0.21%  38 out of 17999

    ✗ create PR
      ↳  99% — ✓ 4503 / ✗ 13
    ✗ set active
      ↳  99% — ✓ 4494 / ✗ 6
    ✗ get team
      ↳  99% — ✓ 4428 / ✗ 11
    ✗ get reviews
      ↳  99% — ✓ 4536 / ✗ 8

    HTTP
    http_req_duration..............: avg=453.6ms  min=2.88ms med=150.27ms max=20.52s p(90)=825.78ms p(95)=1.84s
      { expected_response:true }...: avg=448.81ms min=2.88ms med=149.23ms max=20.52s p(90)=822.46ms p(95)=1.82s
    http_req_failed................: 0.22%  50 out of 22499
    http_reqs......................: 22499  506.682287/s

    EXECUTION
    iteration_duration.............: avg=555.93ms min=8.42ms med=284.96ms max=19.44s p(90)=1.43s    p(95)=2.29s
    iterations.....................: 17999  405.341326/s
    vus............................: 1      min=1           max=499
    vus_max........................: 500    min=500         max=500

    NETWORK
    data_received..................: 127 MB 2.9 MB/s
    data_sent......................: 3.9 MB 88 kB/s

running (0m44.4s), 000/500 VUs, 17999 complete and 0 interrupted iterations
default ✓ [ 100% ] 000/500 VUs  40s

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

Тест #4

нагрузка:
1000 виртуальных пользователей
7 запросов за итерацию без пауз
1000 одновременных соединений

результаты:
RPS: 339 (сервер замедлился но не умер)
ошибки: 12.89%
latency p95: 6.21s
итерации: 42 в секунду

вывод k6:
  █ THRESHOLDS

    http_req_duration
    ✓ 'p(95)<10000' p(95)=6.21s

    http_req_failed
    ✓ 'rate<0.9' rate=12.89%


  █ TOTAL RESULTS

    HTTP
    http_req_duration..............: avg=1.75s  min=1.97ms   med=356.6ms  max=1m2s p(90)=4.16s  p(95)=6.21s
      { expected_response:true }...: avg=1.64s  min=2.17ms   med=343.22ms max=1m1s p(90)=4.18s  p(95)=6.18s
    http_req_failed................: 12.89% 2824 out of 21904
    http_reqs......................: 21904  338.971212/s

    EXECUTION
    iteration_duration.............: avg=13.63s min=381.98ms med=9.93s    max=1m1s p(90)=31.98s p(95)=38.15s
    iterations.....................: 2738   42.371402/s
    vus............................: 8      min=8             max=1000
    vus_max........................: 1000   min=1000          max=1000

    NETWORK
    data_received..................: 181 MB 2.8 MB/s
    data_sent......................: 3.6 MB 55 kB/s

running (1m04.6s), 0000/1000 VUs, 2738 complete and 0 interrupted iterations
default ✓ [ 100% ] 0000/1000 VUs  33s

++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

На основе тестов можно сделать следующие выводы:

Сервер стабильно работает при нагрузке до 500 RPS с минимальным уровнем ошибок (0.22%) и приемлемой задержкой до 500 пользователей.
При пиковой нагрузке в 1000 одновременных пользователей значительно возрастает задержка (до 6.21s) и увеличивается процент ошибок (12.89%),
но продолжается обработка запросов. Оптимальная рабочая нагрузка находится в диапазоне 300-500 RPS,
там сервер сохраняет стабильность при задержке ниже 2 секунд.
