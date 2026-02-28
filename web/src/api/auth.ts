import request from './request'


//注册
export function register(data:{
    username: string;
    email: string;
    password: string;
    nickname?:string;}) {
    return request.post('/auth/register', data)
}


//用户名登录
export function login(data:{account: string; password:string}){
    return request.post('/auth/login', data)
}

//发送验证码
export function sendSmsCode(data:{phone: string; purpose: string}){
    return request.post('/auth/sms/send', data)
}

//验证码登录
export function smsLogin(data:{phone:string, code:string}) {
    return request.post('/auth/sms/login', data)
}