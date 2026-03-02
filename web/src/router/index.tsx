
import { createBrowserRouter, Navigate } from 'react-router-dom';
import Layout from '../components/Loyout';
import Home from '../pages/Home';
import Login from '../pages/Login';
import Register from '../pages/Register';
import Pets from '../pages/Pets';
import PetDetail from '../pages/PetDetail';
import PetDiary from '../pages/PetDiary';
import Rescue from '../pages/Rescue';
import Donation from '../pages/Donation';
import Profile from '../pages/Profile';
import Admin from '../pages/Admin';
import Chat from '../pages/chat';
import Notifications from '../pages/Notifications';
import Favorites from '../pages/Favorites';
import Leaderboard from '../pages/Leaderboard';
import { RequireAuth, RequireRole } from '../components/AuthGuard';

const router = createBrowserRouter([
    {
        path: '/',
        element: <Layout />,
        children: [
        { index: true, element: <Home /> },
        { path: 'pets', element: <Pets /> },
        { path: 'pets/:id', element: <PetDetail /> },
        { path: 'pets/:petId/diary', element: <PetDiary /> },
        { path: 'rescue', element: <Rescue /> },
        { path: 'donation', element: <Donation /> },
        { path: 'leaderboard', element: <Leaderboard /> },
        {
            path: 'profile',
            element: (
                <RequireAuth>
                    <Profile />
                </RequireAuth>
        )},
        {
            path: 'notifications',
            element: (
                <RequireAuth>
                    <Notifications />
                </RequireAuth>
        )},
        {
            path: 'favorites',
            element: (
                <RequireAuth>
                    <Favorites />
                </RequireAuth>
        )},
        {
            path: 'admin',
            element: (
                <RequireRole roles={['admin']}>
                    <Admin />
                </RequireRole>
        )},
        {
            path: 'chat',
            element: (
            <RequireAuth>
                <Chat />
            </RequireAuth>
        )},
        ],
    },
    { path: '/login', element: <Login /> },
    { path: '/register', element: <Register /> },
    { path: '*', element: <Navigate to="/" replace /> },
    ]);

export default router;
