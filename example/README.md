## client ##

```
hey -n 100000000 -c 100 -z 10s -m POST -H "Gateway-Jwt: eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJpY2JjIiwiZGF0YSI6eyJmb3JtYXQiOiJzdHJlYW0iLCJjaGFyc2V0IjoiVVRGLTgiLCJoYXNoX3R5cGUiOiJTSEE1MTIiLCJoYXNoZWQiOiJhOTMwYmUwZmM1ODc5OWY0MmI0NjZhYmE3MGQzZGNjMzE4YzA4ODA3M2JkMjg5Mzk4NmJmOThkNzdkM2QxNzM3NGQ5ZWVhNjgzYmRhMzBiMTI0ZGVlNDRiYmI1MTIwZjI3YTM3NmVlMmQ0MzE2ZmZlM2VhYjI0ZGEwZmI5MzEwOSJ9LCJleHAiOjE1ODIxNjA3MzEsImlhdCI6MTU4MjA3NDMzMSwiaXNzIjoiYXBwaWQwMDAwMSIsImp0aSI6Im1zZ2lkLWFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6IiwibmJmIjoxNTgyMDc0MzMxLCJzdWIiOiIvYXBpL3h4eC95eXkvdjIifQ.ZtZlFkaKHTaoABw4dXU0VlsPsdbzNpNuU77ny-saxG06wb0UItPDD-3-3dgnjMouFXml6vUreXvP0MK8XAfTofCs1cbmwmhyMce_G2Wrs2xjcCzwgCeamPwT-hQTP0N6CZz53W1G0VcyuSq6tjtuxKfah0tKxbdLSgbgSnFceuC2O2BunEu6dHp8BK7nsZAwDwSuhDodLOuo1VjhislgihJUrQAqaN_qTkLqRR3JcGdpGTSbhSJhyukLEzRi8KHUjmgTW4D9Okob_30AFAhvJOtWx7iiYc9wk96wBQ9GczeAMIcW8AJVbIOmAlD9Qu1n42Gdv77WHrrDj0gn0rfAiw" -D ./data/10B http://localhost:9090/api/xxx/yyy/v2
```

## gateway ##

```
nohup ./gateway start &
```

## API server ##

```
nohup gudong start -H='Gateway-Apiresp:{"hashed":"a930be0fc58799f42b466aba70d3dcc318c088073bd2893986bf98d77d3d17374d9eea683bda30b124dee44bbb5120f27a376ee2d4316ffe3eab24da0fb93109","resp":{"return_code":0,"return_msg":"success"}}' --body-file=data/10B &
```