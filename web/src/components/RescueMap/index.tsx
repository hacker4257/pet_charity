import { useEffect, useRef } from 'react';
import { loadAMap } from '../../utils/amap';

interface MapPoint {
    id: number;
    title: string;
    longitude: number;
    latitude: number;
    urgency: string;
    status: string;
}

interface RescueMapProps {
    points: MapPoint[];
    onMarkerClick?: (id: number) => void;
}

const urgencyColorMap: Record<string, string> = {
    critical: '#f5222d',
    high: '#fa8c16',
    medium: '#fadb14',
    low: '#52c41a',
};

const RescueMap = ({ points, onMarkerClick }: RescueMapProps) => {
    const mapRef = useRef<HTMLDivElement>(null);
    const mapObjRef = useRef<any>(null);
    const markersRef = useRef<any[]>([]);

    // 初始化地图
    useEffect(() => {
        let map: any;

        const initMap = async () => {
            const AMap = await loadAMap();

            map = new AMap.Map(mapRef.current, {
                zoom: 12,
                center: [116.397428, 39.90923], // 默认北京，后面会定位
            });

            mapObjRef.current = map;

            // 尝试获取用户位置
            const geolocation = new AMap.Geolocation({
                enableHighAccuracy: true,
                timeout: 5000,
            });
            map.addControl(geolocation);
            geolocation.getCurrentPosition((status: string, result: any) => {
                if (status === 'complete') {
                    map.setCenter([result.position.lng, result.position.lat]);
                }
            });
        };

        initMap();

        return () => {
            // 组件卸载时销毁地图
            if (map) {
                map.destroy();
            }
        };
    }, []);

    // 标记点更新
    useEffect(() => {
        const map = mapObjRef.current;
        if (!map) return;

        // 清除旧标记
        markersRef.current.forEach((m) => map.remove(m));
        markersRef.current = [];

        const AMap = (window as any).AMap;
        if (!AMap) return;

        points.forEach((point) => {
            const marker = new AMap.Marker({
                position: [point.longitude, point.latitude],
                title: point.title,
                content: `<div style="
                    width:16px; height:16px; border-radius:50%;
                    background:${urgencyColorMap[point.urgency] || '#1890ff'};
                    border:2px solid #fff; box-shadow:0 2px 6px rgba(0,0,0,0.3);
                "></div>`,
                offset: new AMap.Pixel(-8, -8),
            });

            // 信息窗体
            const infoWindow = new AMap.InfoWindow({
                content: `
                    <div style="padding:8px;min-width:150px">
                        <strong>${point.title}</strong>
                        <div style="margin-top:4px;color:#666">
                            紧急程度：${point.urgency}
                        </div>
                        <div style="color:#666">状态：${point.status}</div>
                    </div>
                `,
                offset: new AMap.Pixel(0, -16),
            });

            marker.on('click', () => {
                infoWindow.open(map, marker.getPosition());
                onMarkerClick?.(point.id);
            });

            map.add(marker);
            markersRef.current.push(marker);
        });
    }, [points, onMarkerClick]);

    return (
        <div
            ref={mapRef}
            style={{ width: '100%', height: '100%', minHeight: 500, borderRadius: 8}}
        />
    );
};

export default RescueMap;
