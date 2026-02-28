import { Navigate } from 'react-router-dom';
import useAuthStore from '../../store/useAuthStore';

// 需要登录
export const RequireAuth = ({ children }: { children: React.ReactNode }) => {
    const { isLoggedIn } = useAuthStore();

    if (!isLoggedIn) {
        return <Navigate to="/login" replace />;
    }

    return <>{children}</>;
};

// 需要特定角色
export const RequireRole = ({
    children,
    roles,
    }: {
    children: React.ReactNode;
    roles: string[];
    }) => {
    const { isLoggedIn, user } = useAuthStore();

    if (!isLoggedIn) {
        return <Navigate to="/login" replace />;
    }

    if (user && !roles.includes(user.role)) {
        return <Navigate to="/" replace />;
    }

    return <>{children}</>;
};