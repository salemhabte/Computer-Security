import React, { useState, useEffect } from 'react';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Card, CardHeader, CardBody } from '../components/ui/Card';
import { Plus, Trash2, Users, Shield, RefreshCw } from 'lucide-react';
import api from '../lib/axios';

export default function DACManager() {
    const [resources, setResources] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showAddModal, setShowAddModal] = useState(false);
    const [newResource, setNewResource] = useState({ type: 'document', label: 'INTERNAL', name: '' }); // 'name' logic missing in backend? Attributes map can hold name.

    // Access Management State
    const [selectedResource, setSelectedResource] = useState(null);
    const [grantForm, setGrantForm] = useState({ email: '', action: 'read' });
    const [message, setMessage] = useState('');

    const [permissions, setPermissions] = useState([]);
    const [loadingPermissions, setLoadingPermissions] = useState(false);
    const [revokeEmail, setRevokeEmail] = useState('');

    useEffect(() => {
        fetchResources();
    }, []);

    useEffect(() => {
        if (selectedResource) {
            fetchPermissions(selectedResource.ID);
        } else {
            setPermissions([]);
        }
    }, [selectedResource]);

    const fetchPermissions = async (resourceId) => {
        setLoadingPermissions(true);
        try {
            const res = await api.get(`/dac/permissions/${resourceId}`);
            setPermissions(res.data || []);
        } catch (err) {
            // If 403 (not owner), we just show empty or handle gracefully
            console.error("Failed to fetch permissions", err);
            setPermissions([]);
        } finally {
            setLoadingPermissions(false);
        }
    };

    const fetchResources = async () => {
        setLoading(true);
        try {
            const res = await api.get('/dac/resources');
            // Backend returns list of resources.
            // If resources are empty, null is possible if backend returns nil, or empty array.
            setResources(res.data || []);
        } catch (err) {
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleCreateResource = async (e) => {
        e.preventDefault();
        try {
            // Storing 'name' in attributes as backend Resource struct doesn't have Name field explicit
            const attributes = { name: newResource.name };

            await api.post('/dac/resources', {
                type: newResource.type,
                label: newResource.label,
                attributes: attributes
            });
            setShowAddModal(false);
            setNewResource({ type: 'document', label: 'INTERNAL', name: '' });
            fetchResources();
            setMessage('Resource created successfully');
        } catch (err) {
            setMessage('Failed to create resource');
        }
    };

    const handleGrant = async (e) => {
        e.preventDefault();
        if (!selectedResource) return;
        try {
            // actions is array
            await api.post('/dac/grant', {
                resource_id: selectedResource.ID,
                subject: grantForm.email,
                actions: [selectedResource.Type + ':' + grantForm.action] // e.g. "document:read"
            });
            setMessage(`Access granted to ${grantForm.email}`);
            setGrantForm({ email: '', action: 'read' });
            fetchPermissions(selectedResource.ID); // Refresh list
        } catch (err) {
            setMessage(err.response?.data?.error || 'Failed to grant access');
        }
    };

    const handleRevoke = async (email) => {
        if (!selectedResource) return;
        if (!confirm(`Revoke access for ${email}?`)) return;

        try {
            await api.post('/dac/revoke', {
                resource_id: selectedResource.ID,
                subject: email
            });
            setMessage(`Access revoked for ${email}`);
            setRevokeEmail('');
            fetchPermissions(selectedResource.ID); // Refresh list
        } catch (err) {
            setMessage('Failed to revoke access');
        }
    };

    return (
        <div className="p-6 max-w-6xl mx-auto">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold text-gray-800">Resource Management (DAC)</h1>
                <Button onClick={() => setShowAddModal(true)} className="flex items-center gap-2">
                    <Plus size={18} /> New Resource
                </Button>
            </div>

            {message && (
                <div className="bg-blue-50 text-blue-800 p-3 rounded-lg mb-4 flex justify-between">
                    <span>{message}</span>
                    <button onClick={() => setMessage('')}>&times;</button>
                </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Resource List */}
                <div className="md:col-span-1 space-y-4">
                    {loading ? <p>Loading...</p> : resources.length === 0 ? <p className="text-gray-500">No resources found.</p> : (
                        resources.map(res => (
                            <div
                                key={res.ID}
                                onClick={() => setSelectedResource(res)}
                                className={`p-4 rounded-lg border cursor-pointer transition-all hover:shadow-md ${selectedResource?.ID === res.ID ? 'border-blue-500 bg-blue-50 ring-2 ring-blue-200' : 'bg-white border-gray-200'}`}
                            >
                                <div className="flex justify-between items-start">
                                    <div>
                                        <h3 className="font-semibold text-gray-900">{res.Attributes?.name || 'Unnamed Resource'}</h3>
                                        <p className="text-xs text-gray-500 font-mono mt-1">{res.ID.substring(0, 8)}...</p>
                                    </div>
                                    <span className={`text-xs px-2 py-1 rounded-full ${res.Label === 'CONFIDENTIAL' ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'}`}>
                                        {res.Label}
                                    </span>
                                </div>
                                <p className="text-sm text-gray-600 mt-2 capitalize">{res.Type}</p>
                            </div>
                        ))
                    )}
                </div>

                {/* Access Control Panel */}
                <div className="md:col-span-2">
                    {selectedResource ? (
                        <Card className="h-full">
                            <CardHeader>
                                <h2 className="text-xl font-bold flex items-center gap-2">
                                    <Shield className="text-purple-600" />
                                    Access Control: {selectedResource.Attributes?.name}
                                </h2>
                                <p className="text-sm text-gray-500 mt-1">Manage who can access this {selectedResource.Type}</p>
                            </CardHeader>
                            <CardBody>
                                <div className="mb-6 bg-gray-50 p-4 rounded-lg border border-gray-100">
                                    <h3 className="font-semibold mb-3">Grant Access</h3>
                                    <form onSubmit={handleGrant} className="flex gap-3 items-end">
                                        <Input
                                            label="User Email"
                                            placeholder="user@example.com"
                                            value={grantForm.email}
                                            onChange={(e) => setGrantForm({ ...grantForm, email: e.target.value })}
                                            required
                                            className="flex-1"
                                        />
                                        <div className="w-1/3">
                                            <label className="block text-sm font-medium text-gray-700 mb-1">Permission</label>
                                            <select
                                                className="w-full px-3 py-2 border rounded-lg"
                                                value={grantForm.action}
                                                onChange={(e) => setGrantForm({ ...grantForm, action: e.target.value })}
                                            >
                                                <option value="read">Read</option>
                                                <option value="write">Write</option>
                                                <option value="delete">Delete</option>
                                            </select>
                                        </div>
                                        <Button type="submit">Grant</Button>
                                    </form>
                                </div>

                                <div>
                                    <div>
                                        <h3 className="font-semibold mb-3 flex items-center gap-2">
                                            <Users size={18} />
                                            Active Permissions
                                        </h3>

                                        {loadingPermissions ? (
                                            <div className="text-gray-500 text-sm">Loading permissions...</div>
                                        ) : permissions.length === 0 ? (
                                            <div className="text-gray-500 text-sm italic">No active permissions found (private).</div>
                                        ) : (
                                            <div className="space-y-2">
                                                {permissions.map((perm, idx) => (
                                                    <div key={idx} className="flex justify-between items-center bg-white p-2 border rounded-lg shadow-sm">
                                                        <div>
                                                            <p className="font-medium text-sm">{perm.Subject}</p>
                                                            <div className="flex gap-1 mt-1">
                                                                {perm.Actions.map((action, i) => (
                                                                    <span key={i} className="text-xs bg-blue-100 text-blue-700 px-1.5 py-0.5 rounded">
                                                                        {action.split(':')[1] || action}
                                                                    </span>
                                                                ))}
                                                            </div>
                                                        </div>
                                                        <Button
                                                            variant="text"
                                                            className="text-red-600 hover:text-red-800 hover:bg-red-50 p-1 h-auto"
                                                            onClick={() => handleRevoke(perm.Subject)}
                                                        >
                                                            <Trash2 size={16} />
                                                        </Button>
                                                    </div>
                                                ))}
                                            </div>
                                        )}

                                        {/* Manual Revoke Fallback (Optional, but keeping active list is better) */}
                                        <h4 className="text-sm font-bold mt-6 mb-2 text-gray-500 uppercase text-xs">Manual Revoke</h4>
                                        <div className="flex gap-2">
                                            <Input
                                                placeholder="Enter email to revoke"
                                                className="flex-1"
                                                value={revokeEmail}
                                                onChange={(e) => setRevokeEmail(e.target.value)}
                                            />
                                            <Button variant="danger" onClick={() => handleRevoke(revokeEmail)}>Revoke</Button>
                                        </div>
                                    </div>
                            </CardBody>
                        </Card>
                    ) : (
                        <div className="h-full flex items-center justify-center bg-gray-50 rounded-xl border-2 border-dashed border-gray-200 p-12 text-gray-400">
                            <p>Select a resource to manage access</p>
                        </div>
                    )}
                </div>
            </div>

            {/* Add Resource Modal */}
            {showAddModal && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
                    <Card className="w-full max-w-md">
                        <CardHeader>
                            <h2 className="text-xl font-bold">Create New Resource</h2>
                        </CardHeader>
                        <CardBody>
                            <form onSubmit={handleCreateResource} className="space-y-4">
                                <Input
                                    label="Resource Name"
                                    required
                                    value={newResource.name}
                                    onChange={(e) => setNewResource({ ...newResource, name: e.target.value })}
                                />
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Sensitivity Label</label>
                                    <select
                                        className="w-full px-3 py-2 border rounded-lg"
                                        value={newResource.label}
                                        onChange={(e) => setNewResource({ ...newResource, label: e.target.value })}
                                    >
                                        <option value="PUBLIC">PUBLIC</option>
                                        <option value="INTERNAL">INTERNAL</option>
                                        <option value="CONFIDENTIAL">CONFIDENTIAL</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Type</label>
                                    <select
                                        className="w-full px-3 py-2 border rounded-lg"
                                        value={newResource.type}
                                        onChange={(e) => setNewResource({ ...newResource, type: e.target.value })}
                                    >
                                        <option value="document">Document</option>
                                        <option value="record">Record</option>
                                        <option value="file">File</option>
                                    </select>
                                </div>
                                <div className="flex justify-end gap-2 mt-6">
                                    <Button type="button" variant="secondary" onClick={() => setShowAddModal(false)}>Cancel</Button>
                                    <Button type="submit">Create</Button>
                                </div>
                            </form>
                        </CardBody>
                    </Card>
                </div>
            )}
        </div>
    );
}
