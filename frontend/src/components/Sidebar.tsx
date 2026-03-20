"use client";

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { LayoutDashboard, Users, CheckSquare, CalendarDays, Award, LogOut } from 'lucide-react';
import { ProfileAvatar } from '@/components/ui/ProfileAvatar';
import { useRouter } from 'next/navigation';

export const Sidebar = () => {
  const pathname = usePathname();
  const { user, logout } = useAuth();
  const isGuru = user?.role === 'guru';
  const router = useRouter();

  const menuItems = [
    { name: 'Dashboard',    path: '/dashboard',          icon: LayoutDashboard, color: '#818cf8' },
    { name: isGuru ? 'Students' : 'My Profile', path: '/dashboard/students', icon: Users,          color: '#34d399' },
    { name: 'Tugas',        path: '/dashboard/tasks',    icon: CheckSquare,     color: '#f472b6' },
    { name: 'Absensi',      path: '/dashboard/absensi',  icon: CalendarDays,    color: '#fbbf24' },
    { name: 'Nilai',        path: '/dashboard/grades',   icon: Award,           color: '#60a5fa' },
  ];

  const handleLogout = () => { logout(); router.push('/login'); };

  return (
    <aside
      className="flex flex-col h-full relative overflow-hidden"
      style={{
        width: 260,
        minWidth: 260,
        flexShrink: 0,
        background: 'linear-gradient(180deg, #0d1117 0%, #111827 100%)',
        borderRight: '1px solid rgba(255,255,255,0.06)',
      }}
    >
      {/* Ambient glow */}
      <div className="pointer-events-none absolute -top-20 -left-10 w-60 h-60 rounded-full"
        style={{ background: 'radial-gradient(circle, rgba(124,58,237,0.15), transparent 70%)' }} />

      {/* Logo */}
      <div className="relative px-5 py-4" style={{ borderBottom: '1px solid rgba(255,255,255,0.06)' }}>
        <div className="flex items-center gap-3">
          <div className="relative">
            <div className="w-9 h-9 rounded-xl flex items-center justify-center text-lg"
              style={{ background: 'linear-gradient(135deg, #7c3aed, #4f46e5)', boxShadow: '0 4px 16px rgba(124,58,237,0.4)' }}>
              🏫
            </div>
          </div>
          <div>
            <h1 className="text-sm font-bold text-white tracking-tight">SchoolApp</h1>
            <p className="text-[10px] font-medium uppercase tracking-widest" style={{ color: '#6d7280' }}>Portal Manajemen</p>
          </div>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-3 space-y-1 overflow-y-auto">
        <p className="text-[10px] font-bold uppercase tracking-widest px-3 mb-3" style={{ color: '#374151' }}>Menu</p>
        {menuItems.map((item) => {
          const isActive = pathname === item.path || (item.path !== '/dashboard' && pathname.startsWith(item.path));
          const Icon = item.icon;
          return (
            <Link
              key={item.name}
              href={item.path}
              className="relative flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-200 group"
              style={isActive ? {
                background: 'rgba(124,58,237,0.15)',
                border: '1px solid rgba(124,58,237,0.25)',
              } : {
                border: '1px solid transparent',
              }}
            >
              {/* Active accent line */}
              {isActive && (
                <div className="absolute left-0 top-1/2 -translate-y-1/2 w-0.5 h-5 rounded-r-full"
                  style={{ background: item.color }} />
              )}
              {/* Icon */}
              <div className="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 transition-all"
                style={isActive ? {
                  background: `rgba(${hexToRgb(item.color)},0.2)`,
                  boxShadow: `0 2px 10px rgba(${hexToRgb(item.color)},0.25)`,
                } : {
                  background: 'rgba(255,255,255,0.05)',
                }}>
                <Icon size={15} style={{ color: isActive ? item.color : '#6b7280' }} />
              </div>
              <span className="text-sm font-medium" style={{ color: isActive ? '#f9fafb' : '#9ca3af' }}>
                {item.name}
              </span>
              {isActive && (
                <div className="ml-auto w-1.5 h-1.5 rounded-full" style={{ background: item.color }} />
              )}
            </Link>
          );
        })}
      </nav>

      {/* User footer */}
      <div className="p-3" style={{ borderTop: '1px solid rgba(255,255,255,0.06)' }}>
        <div className="flex items-center gap-3 p-3 rounded-xl mb-2"
          style={{ background: 'rgba(255,255,255,0.04)', border: '1px solid rgba(255,255,255,0.07)' }}>
          <ProfileAvatar profileImage={user?.profile_image} username={user?.username} size="md" />
          <div className="flex-1 min-w-0">
            <p className="text-sm font-semibold text-white truncate">{user?.username || 'User'}</p>
            <p className="text-xs" style={{ color: '#6b7280' }}>
              {user?.role === 'guru' ? '👨‍🏫 Guru' : '👨‍🎓 Siswa'}
            </p>
          </div>
        </div>
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-2 px-3 py-2.5 rounded-xl text-sm font-medium transition-all"
          style={{ color: '#6b7280' }}
          onMouseEnter={e => { (e.currentTarget as any).style.background = 'rgba(239,68,68,0.1)'; (e.currentTarget as any).style.color = '#f87171'; }}
          onMouseLeave={e => { (e.currentTarget as any).style.background = 'transparent'; (e.currentTarget as any).style.color = '#6b7280'; }}
        >
          <LogOut size={14} />
          Keluar
        </button>
      </div>
    </aside>
  );
};

// Helper: convert hex color to "r,g,b" string for rgba
function hexToRgb(hex: string): string {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return `${r},${g},${b}`;
}
