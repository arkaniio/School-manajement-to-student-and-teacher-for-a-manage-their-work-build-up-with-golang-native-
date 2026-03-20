"use client";

import React, { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Bell, Camera, X, Sparkles, LogOut } from 'lucide-react';
import api from '@/lib/api';
import { ProfileAvatar } from '@/components/ui/ProfileAvatar';
import { usePathname } from 'next/navigation';

const PAGE_TITLES: Record<string, string> = {
  '/dashboard': 'Dashboard',
  '/dashboard/students': 'Students',
  '/dashboard/tasks': 'Tugas',
  '/dashboard/absensi': 'Absensi',
  '/dashboard/grades': 'Nilai',
};

export const Navbar = () => {
  const { user, updateUser, logout } = useAuth();
  const pathname = usePathname();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState({ username: '', email: '', password: '' });
  const [profileImage, setProfileImage] = useState<File | null>(null);
  const [imagePreview, setImagePreview] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const pageTitle = PAGE_TITLES[pathname] || 'Dashboard';

  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;
    setLoading(true);
    setError('');
    try {
      const data = new FormData();
      if (formData.username) data.append('username', formData.username);
      if (formData.email) data.append('email', formData.email);
      if (formData.password) data.append('password', formData.password);
      if (profileImage) data.append('profile_image', profileImage);

      const res = await api.patch(`/users/${user.id}`, data, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });

      const updated = res.data?.data;
      updateUser({
        username: updated?.Username || updated?.username || formData.username || user.username,
        email: updated?.Email || updated?.email || formData.email || user.email,
        profile_image: updated?.Profile_Image || updated?.profile_image || user.profile_image,
      });
      setIsModalOpen(false);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to update profile.');
    } finally {
      setLoading(false);
    }
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      setProfileImage(file);
      setImagePreview(URL.createObjectURL(file));
    }
  };

  const openProfile = () => {
    setFormData({ username: user?.username || '', email: user?.email || '', password: '' });
    setProfileImage(null);
    setImagePreview(null);
    setError('');
    setIsModalOpen(true);
  };

  return (
    <>
      <header className="h-14 flex items-center justify-between px-6 z-10 sticky top-0"
        style={{
          background: 'rgba(10,14,26,0.85)',
          backdropFilter: 'blur(16px)',
          WebkitBackdropFilter: 'blur(16px)',
          borderBottom: '1px solid rgba(255,255,255,0.06)',
        }}>

        {/* Left: Page title */}
        <div className="flex items-center gap-3">
          <h2 className="text-base font-bold text-white">{pageTitle}</h2>
          <span className="badge" style={user?.role === 'guru'
            ? { background: 'rgba(99,102,241,0.15)', color: '#818cf8', border: '1px solid rgba(99,102,241,0.3)' }
            : { background: 'rgba(16,185,129,0.15)', color: '#34d399', border: '1px solid rgba(16,185,129,0.3)' }}>
            {user?.role === 'guru' ? '👨‍🏫 Guru' : '👨‍🎓 Siswa'}
          </span>
        </div>

        {/* Right: Actions */}
        <div className="flex items-center gap-2">
          {/* Notification bell */}
          <button className="relative p-2 rounded-xl transition-all"
            style={{ color: '#6b7280' }}
            onMouseEnter={e => { (e.currentTarget as any).style.background = 'rgba(255,255,255,0.07)'; (e.currentTarget as any).style.color = '#f9fafb'; }}
            onMouseLeave={e => { (e.currentTarget as any).style.background = 'transparent'; (e.currentTarget as any).style.color = '#6b7280'; }}>
            <Bell size={18} />
            <span className="absolute top-1.5 right-1.5 w-2 h-2 rounded-full animate-pulse"
              style={{ background: '#f43f5e', boxShadow: '0 0 6px #f43f5e' }} />
          </button>

          <div className="w-px h-5 mx-1" style={{ background: 'rgba(255,255,255,0.08)' }} />

          {/* User info + avatar */}
          <div className="flex items-center gap-3 cursor-pointer group" onClick={openProfile}>
            <div className="text-right hidden sm:block">
              <p className="text-sm font-semibold text-white leading-normal transition-colors group-hover:text-violet-300">
                {user?.username}
              </p>
              <p className="text-[11px]" style={{ color: '#6b7280' }}>{user?.email}</p>
            </div>
            <div className="relative">
              <div className="rounded-full ring-2 ring-transparent group-hover:ring-violet-500/50 transition-all"
                style={{ padding: 2 }}>
                <ProfileAvatar profileImage={user?.profile_image} username={user?.username} size="md" />
              </div>
              <div className="absolute -bottom-0.5 -right-0.5 w-3 h-3 rounded-full border-2"
                style={{ background: '#10b981', borderColor: '#0a0e1a' }} />
            </div>
          </div>

          <div className="w-px h-5 mx-1 hidden sm:block" style={{ background: 'rgba(255,255,255,0.08)' }} />

          {/* Logout Button */}
          <button onClick={logout} className="relative p-2 rounded-xl transition-all text-red-400 hover:text-red-300 hover:bg-red-400/10" title="Logout">
            <LogOut size={18} />
          </button>
        </div>
      </header>

      {/* Edit Profile Modal */}
      {isModalOpen && (
        <div className="modal-backdrop" onClick={() => setIsModalOpen(false)}>
          <div className="modal-box" onClick={e => e.stopPropagation()}>
            {/* Header gradient */}
            <div className="h-1 rounded-t-[20px]"
              style={{ background: 'linear-gradient(90deg, #7c3aed, #4f46e5, #06b6d4)' }} />

            <div className="flex items-center justify-between px-5 py-3"
              style={{ borderBottom: '1px solid rgba(255,255,255,0.07)' }}>
              <h3 className="text-base font-bold text-white flex items-center gap-2">
                <Sparkles size={16} style={{ color: '#a78bfa' }} /> Edit Profile
              </h3>
              <button onClick={() => setIsModalOpen(false)}
                className="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-white/10 transition-all">
                <X size={16} />
              </button>
            </div>

            <form onSubmit={handleUpdateProfile} className="p-5 space-y-4">
              {error && (
                <div className="p-3 rounded-xl text-sm" style={{ background: 'rgba(239,68,68,0.1)', border: '1px solid rgba(239,68,68,0.25)', color: '#f87171' }}>
                  ⚠ {error}
                </div>
              )}

              {/* Avatar Upload */}
              <div className="flex justify-center">
                <div className="relative group">
                  <div className="w-20 h-20 rounded-full overflow-hidden"
                    style={{ border: '3px solid rgba(139,92,246,0.5)', boxShadow: '0 0 20px rgba(139,92,246,0.3)' }}>
                    {imagePreview ? (
                      <img src={imagePreview} alt="Preview" className="w-full h-full object-cover" />
                    ) : (
                      <ProfileAvatar profileImage={user?.profile_image} username={user?.username} size="lg" />
                    )}
                  </div>
                  <label className="absolute inset-0 rounded-full bg-black/60 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer">
                    <Camera size={20} className="text-white" />
                    <input type="file" accept="image/png,image/jpeg,image/jpg" className="sr-only" onChange={handleImageChange} />
                  </label>
                </div>
              </div>

              {/* Fields */}
              {[
                { label: 'Username', key: 'username', type: 'text', placeholder: 'Your username' },
                { label: 'Email', key: 'email', type: 'email', placeholder: 'you@example.com' },
                { label: 'Password Baru (opsional)', key: 'password', type: 'password', placeholder: 'Kosongkan jika tidak diganti' },
              ].map(field => (
                <div key={field.key}>
                  <label className="block text-xs font-semibold mb-1.5 uppercase tracking-wider" style={{ color: '#6b7280' }}>
                    {field.label}
                  </label>
                  <input
                    type={field.type}
                    value={(formData as any)[field.key]}
                    onChange={e => setFormData({ ...formData, [field.key]: e.target.value })}
                    placeholder={field.placeholder}
                    className="input-dark"
                  />
                </div>
              ))}

              <div className="flex justify-end gap-3 pt-3" style={{ borderTop: '1px solid rgba(255,255,255,0.07)' }}>
                <button type="button" onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 text-sm font-semibold rounded-xl transition-all"
                  style={{ background: 'rgba(255,255,255,0.06)', color: '#9ca3af', border: '1px solid rgba(255,255,255,0.1)' }}>
                  Batal
                </button>
                <button type="submit" disabled={loading} className="btn-primary" style={{ opacity: loading ? 0.6 : 1 }}>
                  {loading ? 'Menyimpan...' : 'Simpan'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </>
  );
};
