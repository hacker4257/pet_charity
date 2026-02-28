import { useEffect, useState } from 'react';
import { Layout as AntLayout, Menu, Button, Dropdown, Badge } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
    HomeOutlined,
    HeartOutlined,
    AlertOutlined,
    GiftOutlined,
    UserOutlined,
    GlobalOutlined,
    MessageOutlined,
    BellOutlined,
    TrophyOutlined,
} from '@ant-design/icons';
import useAuthStore from '../../store/useAuthStore';
import { getUnreadCount } from '../../api/notification';

const { Header, Content, Footer } = AntLayout;

const Layout = () => {
const navigate = useNavigate();
const location = useLocation();
const { t, i18n } = useTranslation();
const { user, isLoggedIn, logout } = useAuthStore();
const [unreadCount, setUnreadCount] = useState(0);

// 轮询未读通知数
useEffect(() => {
    if (!isLoggedIn) {
        setUnreadCount(0);
        return;
    }
    const fetchUnread = () => {
        getUnreadCount()
            .then((res: any) => setUnreadCount(res.data?.count || 0))
            .catch(() => {});
    };
    fetchUnread();
    const timer = setInterval(fetchUnread, 30000);
    return () => clearInterval(timer);
}, [isLoggedIn]);

const menuItems = [
    { key: '/', icon: <HomeOutlined />, label: t('common.home') },
    { key: '/pets', icon: <HeartOutlined />, label: t('nav.pets') },
    { key: '/rescue', icon: <AlertOutlined />, label: t('nav.rescue') },
    { key: '/donation', icon: <GiftOutlined />, label: t('nav.donation') },
    { key: '/leaderboard', icon: <TrophyOutlined />, label: t('nav.leaderboard') || '排行榜' },
    { key: '/chat', icon: <MessageOutlined />, label: t('nav.chat') || '消息' },
];

// 语言切换
const toggleLanguage = () => {
    const newLang = i18n.language === 'zh-CN' ? 'en-US' : 'zh-CN';
    i18n.changeLanguage(newLang);
};

return (
    <AntLayout style={{ minHeight: '100vh' }}>
    <Header style={{ display: 'flex', alignItems: 'center' }}>
        <div
            style={{ color: '#fff', fontSize: 18, fontWeight: 'bold', marginRight: 40, cursor: 'pointer' }}
            onClick={() => navigate('/')}
        >
        Pet Charity
        </div>

        <Menu
            theme="dark"
            mode="horizontal"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({key}:{key:string}) => navigate(key)}
            style={{ flex: 1 }}
        />

        <Button
        icon={<GlobalOutlined />}
        type="text"
        style={{ color: '#fff', marginRight: 8 }}
        onClick={toggleLanguage}
        >
        {i18n.language === 'zh-CN' ? 'EN' : '中文'}
        </Button>

        {isLoggedIn && (
        <>
            <Button
                type="text"
                style={{ color: '#fff', marginRight: 8 }}
                onClick={() => navigate('/favorites')}
                icon={<HeartOutlined />}
            />
            <Badge count={unreadCount} size="small" offset={[-4, 4]}>
                <Button
                    type="text"
                    style={{ color: '#fff', marginRight: 8 }}
                    onClick={() => navigate('/notifications')}
                    icon={<BellOutlined />}
                />
            </Badge>
        </>
        )}

        {isLoggedIn ? (
        <Dropdown
            menu={{
                items: [
                    { key: 'profile', icon: <UserOutlined />, label: t('nav.profile') },
                    ...(user?.role === 'admin'
                    ? [{ key: 'admin', label: t('nav.admin') }]
                    : []),
                    { type: 'divider' as const },
                    { key: 'logout', label: t('common.logout'), danger: true },
                ],
                onClick: ({key }:{key:string}) => {
                    if (key === 'logout') logout();
                    else navigate(`/${key}`);
                },
            }}
        >
            <Button type="text" style={{ color: '#fff' }}>
            <UserOutlined /> {user?.nickname || user?.username || 'User'}
            </Button>
        </Dropdown>
        ) : (
        <Button type="primary" onClick={() => navigate('/login')}>
            {t('common.login')}
        </Button>
        )}
    </Header>

    <Content style={{ padding: '24px 48px' }}>
        <Outlet />
    </Content>

    <Footer style={{ textAlign: 'center' }}>
        Pet Charity Platform ©2026
    </Footer>
    </AntLayout>
);
};

export default Layout;
