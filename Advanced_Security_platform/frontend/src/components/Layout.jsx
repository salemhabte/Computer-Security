import React from 'react';
import { useAuth } from '../context/AuthContext';
import { Link, Outlet, useLocation } from 'react-router-dom';
import { LayoutDashboard, Lock, User, LogOut, Menu, X } from 'lucide-react';
import { useState } from 'react';

export default function Layout() {
    const { user, logout } = useAuth();
    const location = useLocation();
    const [sidebarOpen, setSidebarOpen] = useState(false);

    const navigation = [
        { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
        { name: 'DAC Manager', href: '/dac', icon: Lock },
        { name: 'Profile', href: '/profile', icon: User },
    ];

    return (
        <div className="min-h-screen bg-gray-50 flex">
            {/* Mobile Sidebar Overlay */}
            {sidebarOpen && (
                <div
                    className="fixed inset-0 bg-black bg-opacity-50 z-20 md:hidden"
                    onClick={() => setSidebarOpen(false)}
                ></div>
            )}

            {/* Sidebar */}
            <aside className={`
        fixed inset-y-0 left-0 w-64 bg-white border-r border-gray-200 transform transition-transform duration-200 ease-in-out z-30
        ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
        md:translate-x-0 md:static
      `}>
                <div className="h-16 flex items-center justify-between px-6 border-b border-gray-100">
                    <span className="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">
                        SecurePlatform
                    </span>
                    <button className="md:hidden" onClick={() => setSidebarOpen(false)}>
                        <X size={20} />
                    </button>
                </div>

                <nav className="p-4 space-y-1">
                    {navigation.map((item) => {
                        const isActive = location.pathname === item.href;
                        return (
                            <Link
                                key={item.name}
                                to={item.href}
                                onClick={() => setSidebarOpen(false)}
                                className={`flex items-center px-4 py-3 rounded-lg transition-colors ${isActive
                                        ? 'bg-blue-50 text-blue-700'
                                        : 'text-gray-600 hover:bg-gray-50'
                                    }`}
                            >
                                <item.icon size={20} className="mr-3" />
                                <span className="font-medium">{item.name}</span>
                            </Link>
                        );
                    })}
                </nav>

                <div className="absolute bottom-0 left-0 right-0 p-4 border-t border-gray-100">
                    <div className="flex items-center gap-3 px-4 py-3 mb-2">
                        <div className="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 font-bold">
                            {user?.email?.[0].toUpperCase()}
                        </div>
                        <div className="overflow-hidden">
                            <p className="text-sm font-medium text-gray-900 truncate">{user?.email}</p>
                            <p className="text-xs text-gray-500">{user?.role}</p>
                        </div>
                    </div>
                    <button
                        onClick={logout}
                        className="w-full flex items-center px-4 py-2 text-sm text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    >
                        <LogOut size={18} className="mr-3" />
                        Sign Out
                    </button>
                </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1 flex flex-col overflow-hidden">
                {/* Mobile Header */}
                <header className="h-16 bg-white border-b border-gray-200 flex items-center px-4 md:hidden">
                    <button onClick={() => setSidebarOpen(true)}>
                        <Menu size={24} className="text-gray-600" />
                    </button>
                    <span className="ml-4 font-bold text-gray-900">SecurePlatform</span>
                </header>

                <div className="flex-1 overflow-auto p-4 md:p-8">
                    <Outlet />
                </div>
            </main>
        </div>
    );
}
