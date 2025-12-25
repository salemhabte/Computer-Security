import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card, CardHeader, CardBody } from '../components/ui/Card';
import { User, Save, LockKeyhole } from 'lucide-react';
import api from '../lib/axios';

export default function Profile() {
    const { user } = useAuth();
    // We should fetch full profile details on load.
    // For now, using what we have and placeholders.

    const [formData, setFormData] = useState({
        bio: 'Security enthusiast.',
        phone: '',
        telegram: '',
    });
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');

    const handleUpdate = async (e) => {
        e.preventDefault();
        setMessage('');

        // Validate Phone Number
        // Validate Phone Number
        const phone = formData.phone;
        const startsWithPlus = /^\+[0-9]{12}$/; // + followed by 12 digits (total 13)
        const startsWithZero = /^0[0-9]{9}$/;   // 0 followed by 9 digits (total 10)

        if (phone) {
            if (phone.startsWith('+')) {
                if (!startsWithPlus.test(phone)) {
                    setMessage('Invalid international format. Must start with + and have 13 characters (e.g., +251911234567)');
                    return;
                }
            } else if (phone.startsWith('0')) {
                if (!startsWithZero.test(phone)) {
                    setMessage('Invalid local format. Must start with 0 and have 10 digits (e.g., 0911234567)');
                    return;
                }
            } else {
                setMessage('Phone number must start with + (for international) or 0 (for local)');
                return;
            }
        }

        setLoading(true);
        try {
            await api.put('/user/edit_profile', {
                personal_bio: formData.bio,
                phone_num: formData.phone,
                telegram_handle: formData.telegram,
            });
            setMessage('Profile updated successfully');
        } catch (err) {
            setMessage('Failed to update profile');
        } finally {
            setLoading(false);
        }
    };

    const [passwordData, setPasswordData] = useState({
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
    });
    const [pwdLoading, setPwdLoading] = useState(false);
    const [pwdMessage, setPwdMessage] = useState('');

    const handleChangePassword = async (e) => {
        e.preventDefault();
        setPwdMessage('');

        if (passwordData.newPassword !== passwordData.confirmPassword) {
            setPwdMessage('New passwords do not match');
            return;
        }

        // Basic frontend check for password strength (optional, backend does it too)
        const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#~$%^&*()_+|<>?:{}]).{8,}$/;
        if (!passwordRegex.test(passwordData.newPassword)) {
            setPwdMessage('New password must contain uppercase, lowercase, number, and special character.');
            return;
        }

        setPwdLoading(true);
        try {
            await api.post('/user/change_password', {
                old_password: passwordData.oldPassword,
                new_password: passwordData.newPassword
            });
            setPwdMessage('Password changed successfully');
            setPasswordData({ oldPassword: '', newPassword: '', confirmPassword: '' });
        } catch (err) {
            setPwdMessage(err.response?.data?.error || 'Failed to change password');
        } finally {
            setPwdLoading(false);
        }
    };

    return (
        <div className="max-w-2xl mx-auto">
            <Card>
                <CardHeader>
                    <h1 className="text-2xl font-bold flex items-center gap-2">
                        <User className="text-blue-600" />
                        Edit Profile
                    </h1>
                </CardHeader>
                <CardBody>
                    <div className="mb-6 flex items-center gap-4">
                        <div className="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center text-2xl font-bold text-blue-600">
                            {user?.email?.[0].toUpperCase()}
                        </div>
                        <div>
                            <p className="font-medium text-lg">{user?.email}</p>
                            <p className="text-gray-500 text-sm capitalize">{user?.role?.toLowerCase()}</p>
                        </div>
                    </div>

                    <form onSubmit={handleUpdate} className="space-y-4">
                        <div className="w-full">
                            <label className="block text-sm font-medium text-gray-700 mb-1">Personal Bio</label>
                            <textarea
                                className="w-full px-3 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                                rows={4}
                                value={formData.bio}
                                onChange={(e) => setFormData({ ...formData, bio: e.target.value })}
                            ></textarea>
                        </div>

                        <Input
                            label="Phone Number"
                            value={formData.phone}
                            onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                            placeholder="+1234567890"
                        />

                        <Input
                            label="Telegram Handle"
                            value={formData.telegram}
                            onChange={(e) => setFormData({ ...formData, telegram: e.target.value })}
                            placeholder="@username"
                        />

                        {message && (
                            <div className={`p-3 rounded-lg text-sm ${message.includes('success') ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
                                {message}
                            </div>
                        )}

                        <div className="flex justify-end pt-4">
                            <Button type="submit" disabled={loading} className="flex items-center gap-2">
                                <Save size={18} />
                                Save Changes
                            </Button>
                        </div>
                    </form>
                </CardBody>
            </Card>

            {/* Security Section */}
            <Card className="mt-8">
                <CardHeader>
                    <h2 className="text-xl font-bold flex items-center gap-2">
                        <LockKeyhole className="text-blue-600" />
                        Security
                    </h2>
                </CardHeader>
                <CardBody>
                    <form onSubmit={handleChangePassword} className="space-y-4">
                        <Input
                            label="Current Password"
                            type="password"
                            required
                            value={passwordData.oldPassword}
                            onChange={(e) => setPasswordData({ ...passwordData, oldPassword: e.target.value })}
                            placeholder="••••••••"
                        />
                        <Input
                            label="New Password"
                            type="password"
                            required
                            value={passwordData.newPassword}
                            onChange={(e) => setPasswordData({ ...passwordData, newPassword: e.target.value })}
                            placeholder="New strong password"
                        />
                        <p className="text-xs text-gray-400">Must contain uppercase, lowercase, number, and special char.</p>
                        <Input
                            label="Confirm New Password"
                            type="password"
                            required
                            value={passwordData.confirmPassword}
                            onChange={(e) => setPasswordData({ ...passwordData, confirmPassword: e.target.value })}
                            placeholder="Confirm new password"
                        />

                        {pwdMessage && (
                            <div className={`p-3 rounded-lg text-sm ${pwdMessage.includes('successfully') ? 'bg-green-50 text-green-700' : 'bg-red-50 text-red-700'}`}>
                                {pwdMessage}
                            </div>
                        )}

                        <div className="flex justify-end pt-2">
                            <Button type="submit" disabled={pwdLoading} className="bg-gray-800 hover:bg-gray-700">
                                {pwdLoading ? 'Updating...' : 'Change Password'}
                            </Button>
                        </div>
                    </form>
                </CardBody>
            </Card>
        </div>
    );
}
