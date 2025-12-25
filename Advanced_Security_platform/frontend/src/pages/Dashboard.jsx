import React from 'react';
import { useAuth } from '../context/AuthContext';
import { Button } from '../components/ui/Button';
import { Card, CardBody } from '../components/ui/Card';
import { Shield, FileText, Settings, Key, User } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function Dashboard() {
    const { user, logout } = useAuth();
    // user object has { email, role, id } based on AuthContext

    return (
        <div className="p-6 max-w-7xl mx-auto">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
                <p className="text-gray-500 mt-2">Welcome back, <span className="font-semibold text-blue-600">{user?.email}</span></p>
                <div className="flex items-center gap-2 mt-2">
                    <span className="px-2 py-1 bg-gray-200 text-gray-700 text-xs rounded-md font-mono">{user?.role}</span>
                </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">

                {/* Quick Actions / DAC */}
                <Card className="hover:shadow-xl transition-shadow border-t-4 border-blue-500">
                    <CardBody>
                        <div className="flex items-center gap-4 mb-4">
                            <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center text-blue-600">
                                <FileText size={24} />
                            </div>
                            <div>
                                <h3 className="font-bold text-lg">My Resources</h3>
                                <p className="text-xs text-gray-500">DAC & Access Management</p>
                            </div>
                        </div>
                        <p className="text-sm text-gray-600 mb-4">
                            Create entries, manage documents, and control who can view your sensitive data.
                        </p>
                        <Link to="/dac">
                            <Button variant="primary" className="w-full">Manage Resources</Button>
                        </Link>
                    </CardBody>
                </Card>

                {/* Profile */}
                <Card className="hover:shadow-xl transition-shadow border-t-4 border-purple-500">
                    <CardBody>
                        <div className="flex items-center gap-4 mb-4">
                            <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center text-purple-600">
                                <User size={24} />
                            </div>
                            <div>
                                <h3 className="font-bold text-lg">Profile</h3>
                                <p className="text-xs text-gray-500">Personal Info</p>
                            </div>
                        </div>
                        <p className="text-sm text-gray-600 mb-4">
                            Update your bio, contact details, and view your clearance status.
                        </p>
                        <Link to="/profile">
                            <Button variant="secondary" className="w-full">Edit Profile</Button>
                        </Link>
                    </CardBody>
                </Card>

                {/* Security / Admin */}
                {user?.role !== 'USER' && (
                    <Card className="hover:shadow-xl transition-shadow border-t-4 border-red-500">
                        <CardBody>
                            <div className="flex items-center gap-4 mb-4">
                                <div className="w-12 h-12 bg-red-100 rounded-lg flex items-center justify-center text-red-600">
                                    <Shield size={24} />
                                </div>
                                <div>
                                    <h3 className="font-bold text-lg">Admin Console</h3>
                                    <p className="text-xs text-gray-500">System Management</p>
                                </div>
                            </div>
                            <p className="text-sm text-gray-600 mb-4">
                                Manage roles, view audit logs, and trigger backups.
                            </p>
                            <Link to="/admin">
                                <Button variant="outline" className="w-full border-red-200 text-red-600 hover:bg-red-50">Open Console</Button>
                            </Link>
                        </CardBody>
                    </Card>
                )}
            </div>
        </div>
    );
}
