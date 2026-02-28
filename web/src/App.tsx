import { useEffect } from 'react';
import { RouterProvider } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import router from './router';
import useAuthStore from './store/useAuthStore';

const App = () => {
  const { isLoggedIn, fetchUser } = useAuthStore();

  // 页面加载时，如果有 token 就拉取用户信息
  useEffect(() => {
    if (isLoggedIn) {
      fetchUser();
    }
  }, []);

  return (
    <ConfigProvider locale={zhCN}>
      <RouterProvider router={router} />
    </ConfigProvider>
  );
};

export default App;
