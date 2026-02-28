import { create } from 'zustand';
import { getMe } from '../api/user';

interface User {
    id: number;
    username: string;
    email: string;
    phone: string;
    nickname: string;
    avatar: string;
    role: string;
    language: string;
}


interface AuthState {
    user: User | null;
    isLoggedIn: boolean;

    //
    loginSuccess: (accessToken: string, refreshToken: string) => Promise<void>;
    logout: ()=> void;

    fetchUser: ()=>Promise<void>;
}


const useAuthStore = create<AuthState>((set) =>({
    user: null,
    isLoggedIn: !!localStorage.getItem('access_token'),

    loginSuccess: async (accessToken, refreshToken) => {
        localStorage.setItem('access_token', accessToken);
        localStorage.setItem('refresh_token', refreshToken);
        set({ isLoggedIn: true});

        try {
            const res: any = await getMe();
            set({user: res.data});
        }catch{
            //
        }
    },
    logout: ()=> {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        set({user:null, isLoggedIn: false});
        window.location.href = '/login';
    },
    fetchUser: async () => {
        try{
            const res: any = await getMe();
            set({user:res.data, isLoggedIn:true});
        }catch{
            set({user:null, isLoggedIn:false})
        }
    },
}));

export default useAuthStore