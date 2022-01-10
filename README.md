# Endpoints

### _All Symbol should be UPPERCASE and follow by USDT, BUSD_

1. Get All Data include prices by this endpoint  
   request  

   ```
   GET: {{host}}/api/coins?symbols=SOLUSDT, BTCBUSD
   ```  

   response  

   ```  
   [{
       "symbol": string,
        "closePrices": float64[],
        "svg": string
   }]
   ```  

2. Get Svg  
   request

   ```
    GET: {{host}}/api/svg/SOLUSDT

   ```

    response example

    ```

    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="720" height="250">\n<path  d="M 0 0 L 720 0 L 720 250 L 0 250 L 0 0" style="stroke-width:0;stroke:rgba(0,0,0,0.0);fill:rgba(0,0,0,0.0)"/> <path  d="M 5 5 L 715 5 L 715 245 L 5 245 L 5 5" style="stroke-width:0;stroke:rgba(0,0,0,0.0);fill:rgba(0,0,0,0.0)"/> <path  d="M 5 191 L 9 199 L 12 208 L 15 216 L 18 223 L 21 230 L 24 235 L 27 240 L 30 243 L 33 245 L 36 244 L 39 242 L 43 239 L 46 234 L 49 228 L 52 222 L 55 216 L 58 210 L 61 205 L 64 200 L 67 197 L 70 195 L 73 194 L 76 195 L 80 195 L 83 196 L 86 196 L 89 196 L 92 195 L 95 192 L 98 187 L 101 179 L 104 170 L 107 159 L 110 147 L 114 134 L 117 121 L 120 108 L 123 94 L 126 82 L 129 70 L 132 60 L 135 51 L 138 43 L 141 37 L 144 33 L 147 30 L 151 28 L 154 28 L 157 29 L 160 32 L 163 36 L 166 42 L 169 48 L 172 55 L 175 61 L 178 67 L 181 71 L 185 74 L 188 76 L 191 75 L 194 71 L 197 65 L 200 58 L 203 49 L 206 40 L 209 31 L 212 22 L 215 15 L 218 9 L 222 5 L 225 5 L 228 6 L 231 10 L 234 15 L 237 22 L 240 29 L 243 37 L 246 46 L 249 55 L 252 63 L 256 70 L 259 76 L 262 82 L 265 88 L 268 93 L 271 97 L 274 101 L 277 105 L 280 108 L 283 111 L 286 114 L 289 118 L 293 121 L 296 124 L 299 126 L 302 129 L 305 132 L 308 135 L 311 137 L 314 140 L 317 142 L 320 145 L 323 147 L 327 150 L 330 152 L 333 154 L 336 156 L 339 158 L 342 160 L 345 161 L 348 162 L 351 164 L 354 165 L 357 166 L 360 167 L 364 168 L 367 169 L 370 170 L 373 171 L 376 173 L 379 174 L 382 176 L 385 178 L 388 179 L 391 180 L 394 181 L 398 182 L 401 182 L 404 181 L 407 180 L 410 178 L 413 175 L 416 172 L 419 168 L 422 165 L 425 162 L 428 159 L 431 157 L 435 156 L 438 156 L 441 157 L 444 159 L 447 162 L 450 165 L 453 168 L 456 170 L 459 172 L 462 172 L 465 171 L 469 167 L 472 162 L 475 155 L 478 146 L 481 136 L 484 126 L 487 115 L 490 105 L 493 95 L 496 87 L 499 80 L 502 75 L 506 71 L 509 69 L 512 69 L 515 70 L 518 72 L 521 75 L 524 78 L 527 82 L 530 86 L 533 90 L 536 94 L 540 98 L 543 102 L 546 105 L 549 107 L 552 109 L 555 110 L 558 110 L 561 110 L 564 108 L 567 105 L 570 101 L 573 98 L 577 94 L 580 90 L 583 88 L 586 86 L 589 85 L 592 85 L 595 87 L 598 91 L 601 96 L 604 101 L 607 107 L 611 114 L 614 120 L 617 126 L 620 132 L 623 136 L 626 139 L 629 142 L 632 143 L 635 144 L 638 144 L 641 144 L 644 144 L 648 144 L 651 144 L 654 144 L 657 144 L 660 145 L 663 147 L 666 148 L 669 150 L 672 152 L 675 154 L 678 156 L 682 158 L 685 159 L 688 161 L 691 162 L 694 163 L 697 164 L 700 165 L 703 165 L 706 166 L 709 167 L 712 167 L 715 167" style="stroke-width:5;stroke:rgba(50,161,128,1.0);fill:none"/> </svg> 
    
    ```  

3. fetch kline hourly will start fetching on server started buy also can start by API
   request

   ```
    POST: {{host}}/routine/start
   ```

4. stop fetching kline
   request

   ```
    POST: {{host}}/routine/stop
   ```

5. force update all kline that were fetched in the past
    request

    ```
    POST: {{host}}/routine/force-update
    ```

#### if data do not exist it will fetch from binance and store in redis first then response to you next  

#### all kline will store in Redis and it will expire in a hour
