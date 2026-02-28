import AMapLoader from '@amap/amap-jsapi-loader';

// 这里填你自己的 key，后面可以改成环境变量
const AMAP_KEY = 'your_amap_key_here';
const SECURITY_CODE = 'your_security_code_here';

let mapInstance: any = null;

export async function loadAMap(): Promise<any> {
    if (mapInstance) return mapInstance;

    // 高德安全密钥（2021年12月后申请的 key 必须配置）
    (window as any)._AMapSecurityConfig = {
        securityJsCode: SECURITY_CODE,
    };

    mapInstance = await AMapLoader.load({
        key: AMAP_KEY,
        version: '2.0',
        plugins: ['AMap.Marker', 'AMap.InfoWindow', 'AMap.Geolocation'],
    });

    return mapInstance;
}