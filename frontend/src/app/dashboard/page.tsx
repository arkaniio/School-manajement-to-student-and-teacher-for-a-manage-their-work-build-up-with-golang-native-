"use client";

import React, { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { BookOpen, CalendarCheck, GraduationCap, CheckCircle, Clock, Award, TrendingUp, Zap, Star } from 'lucide-react';
import api from '@/lib/api';
import Link from 'next/link';

interface StudentData {
  id: string;
  full_name: string;
  kelas: string;
  jurusan: string;
  absen: number;
  wali_kelas: string;
  mapel_students: string;
  student_profile: string;
}

export default function DashboardPage() {
  const { user } = useAuth();
  const isGuru = user?.role === 'guru';
  const isSiswa = user?.role === 'siswa';
  const [tasksCount, setTasksCount] = useState(0);
  const [studentData, setStudentData] = useState<StudentData | null>(null);
  const [tasksData, setTasksData] = useState<any[]>([]);
  const [gradesData, setGradesData] = useState<any[]>([]);
  const [loadingStats, setLoadingStats] = useState(true);

  useEffect(() => {
    const controller = new AbortController();
    const fetchData = async () => {
      setLoadingStats(true);
      try {
        if (isSiswa) {
          const res = await api.get('/students', { signal: controller.signal });
          const data = res.data?.data;
          if (data) setStudentData(data);
        }
        if (isGuru) {
          const [tasksRes, gradesRes] = await Promise.all([
            api.get('/tasks', { signal: controller.signal }),
            api.get('/grades', { signal: controller.signal })
          ]);
          const tasks = tasksRes.data?.data || [];
          setTasksData(tasks);
          setTasksCount(Array.isArray(tasks) ? tasks.length : 0);
          setGradesData(gradesRes.data?.data || []);
        }
      } catch (err: any) {
        if (err?.code === 'ERR_CANCELED' || err?.name === 'CanceledError' || err?.name === 'AbortError') return;
      } finally {
        setLoadingStats(false);
      }
    };
    fetchData();
    return () => controller.abort();
  }, [isSiswa, isGuru]);

  const getGreeting = () => {
    const h = new Date().getHours();
    if (h < 12) return { text: 'Selamat Pagi', emoji: '☀️' };
    if (h < 17) return { text: 'Selamat Siang', emoji: '🌤️' };
    return { text: 'Selamat Sore', emoji: '🌙' };
  };
  const greeting = getGreeting();

  const today = new Date().toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });

  // ── GURU DASHBOARD ──────────────────────────────────────
  if (isGuru) {
    const uniqueStudents = new Set(tasksData.map(t => t.Student_Id).filter(Boolean)).size;

    return (
      <div className="space-y-8 animate-fade-in-up">

        {/* Hero Banner */}
        <div className="relative overflow-hidden rounded-3xl p-8 md:p-10"
          style={{
            background: 'linear-gradient(135deg, #1e1333 0%, #2d1b69 50%, #1e1b4b 100%)',
            border: '1px solid rgba(139,92,246,0.25)',
            boxShadow: '0 20px 60px rgba(0,0,0,0.4), inset 0 1px 0 rgba(255,255,255,0.05)',
          }}>
          {/* Decoration orbs */}
          <div className="absolute top-0 right-0 w-80 h-80 rounded-full pointer-events-none"
            style={{ background: 'radial-gradient(circle, rgba(139,92,246,0.25), transparent 65%)', transform: 'translate(30%,-30%)' }} />
          <div className="absolute inset-0 opacity-5 pointer-events-none"
            style={{ backgroundImage: 'radial-gradient(circle, white 1px, transparent 1px)', backgroundSize: '32px 32px' }} />

          <div className="relative z-10 flex items-center gap-6 md:gap-8">
            <div className="relative shrink-0">
              <div className="w-16 h-16 md:w-20 md:h-20 rounded-2xl overflow-hidden flex items-center justify-center text-2xl md:text-3xl font-black text-white shadow-2xl"
                style={{ background: 'linear-gradient(135deg, rgba(139,92,246,0.6), rgba(79,70,229,0.6))', border: '2px solid rgba(255,255,255,0.15)' }}>
                {user?.username?.charAt(0).toUpperCase() || 'G'}
              </div>
              <div className="absolute -bottom-1 -right-1 w-6 h-6 rounded-full border-2 flex items-center justify-center shadow-lg"
                style={{ background: '#10b981', borderColor: '#1e1333', fontSize: 11 }}>✓</div>
            </div>
            <div className="space-y-1">
              <p className="text-xs md:text-sm font-bold tracking-wide leading-none" style={{ color: '#a78bfa' }}>{greeting.emoji} {greeting.text},</p>
              <h1 className="text-2xl md:text-4xl font-black text-white leading-tight">{user?.username || 'Guru'}</h1>
              <div className="flex items-center gap-3 pt-2 flex-wrap">
                <span className="badge badge-purple px-4 py-1 font-bold">👨‍🏫 Guru</span>
                <span className="badge px-4 py-1 font-bold" style={{ background: 'rgba(255,255,255,0.08)', color: '#94a3b8', border: '1px solid rgba(255,255,255,0.12)' }}>
                  📅 {today}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-6 md:gap-8">
          {[
            {
              label: 'Total Tugas', value: loadingStats ? '...' : tasksCount, icon: BookOpen,
              color: '#818cf8', glow: 'rgba(129,140,248,0.2)', desc: 'Aktif',
            },
            {
              label: 'Siswa Aktif', value: loadingStats ? '...' : uniqueStudents, icon: GraduationCap,
              color: '#34d399', glow: 'rgba(52,211,153,0.2)', desc: 'Siswa upload tugas',
            },
            {
              label: 'Nilai Diberikan', value: loadingStats ? '...' : gradesData.length, icon: Award,
              color: '#f472b6', glow: 'rgba(244,114,182,0.2)', desc: 'Tugas dinilai',
            },
            {
              label: 'Hari Ini', value: new Date().toLocaleDateString('id-ID', { day: 'numeric', month: 'short' }), icon: CalendarCheck,
              color: '#fbbf24', glow: 'rgba(251,191,36,0.2)', desc: new Date().toLocaleDateString('id-ID', { weekday: 'short' }),
            },
          ].map((s, i) => (
            <div key={i} className="stat-card glass-hover px-6 py-6 space-y-4" style={{ cursor: 'default' }}>
              <div className="flex items-start justify-between">
                <div className="w-10 h-10 md:w-12 md:h-12 rounded-2xl flex items-center justify-center shadow-lg"
                  style={{ background: s.glow, boxShadow: `0 8px 24px ${s.glow}` }}>
                  <s.icon size={20} className="md:w-6 md:h-6" style={{ color: s.color }} />
                </div>
              </div>
              <div className="space-y-1">
                <p className="text-2xl md:text-3xl font-black text-white leading-tight">{s.value}</p>
                <p className="text-[10px] md:text-xs font-bold uppercase tracking-widest leading-none text-slate-500">{s.label}</p>
              </div>
            </div>
          ))}
        </div>

        {/* Recent Tasks */}
        <div className="glass rounded-3xl overflow-hidden shadow-2xl border border-white/5">
          <div className="flex items-center justify-between px-6 py-5 md:px-8 md:py-6" style={{ borderBottom: '1px solid rgba(255,255,255,0.06)', background: 'rgba(255,255,255,0.01)' }}>
            <h3 className="font-black text-white text-lg flex items-center gap-3 leading-tight">
              <BookOpen size={20} style={{ color: '#8cf881' }} /> Tugas Siswa Terbaru
            </h3>
            <Link href="/dashboard/tasks" className="text-xs font-black transition-colors text-indigo-400 hover:text-indigo-300 uppercase tracking-widest">
              Lihat Lengkap →
            </Link>
          </div>
          {loadingStats ? (
            <div className="flex flex-col gap-4 p-6 md:p-8">
              {[1, 2, 3].map(i => (
                <div key={i} className="h-14 md:h-16 rounded-2xl bg-white/5 animate-pulse" />
              ))}
            </div>
          ) : tasksData.length === 0 ? (
            <div className="text-center py-16 md:py-24 space-y-4">
              <div className="w-16 h-16 md:w-20 md:h-20 rounded-3xl flex items-center justify-center mx-auto shadow-xl"
                style={{ background: 'rgba(129,140,248,0.1)', border: '1px solid rgba(129,140,248,0.2)' }}>
                <BookOpen size={32} style={{ color: '#818cf8' }} />
              </div>
              <div className="space-y-2">
                <p className="font-black text-white text-lg leading-tight">Belum ada tugas siswa</p>
                <p className="text-xs md:text-sm leading-relaxed text-slate-500 max-w-xs mx-auto">Siswa belum mulai mengumpulkan tugas untuk kelas Anda.</p>
              </div>
            </div>
          ) : (
            <div className="p-6 md:p-8 space-y-3">
              {tasksData.slice(0, 5).map((task: any, i: number) => (
                <div key={i} className="flex items-center gap-4 md:gap-6 px-5 py-4 md:py-5 rounded-2xl transition-all bg-white/[.01] border border-white/5 hover:bg-white/[.03] hover:border-white/10 hover:shadow-lg group">
                  <div className="w-8 h-8 md:w-10 md:h-10 rounded-xl flex items-center justify-center shrink-0 text-xs font-black text-white shadow-inner"
                    style={{ background: `hsl(${(i * 60 + 240) % 360},70%,50%)`, opacity: 0.85 }}>{i + 1}</div>
                  <div className="flex-1 min-w-0 space-y-1">
                    <p className="text-sm md:text-base font-bold text-white truncate leading-snug">{task.Name_Task || 'Untitled'}</p>
                    <p className="text-[11px] md:text-xs truncate text-slate-500 leading-none">👤 {task.Students?.Full_Name || 'Unknown Student'} • {task.MapelTask || 'Umum'}</p>
                  </div>
                  <span className="text-[11px] md:text-xs shrink-0 badge badge-info px-4 py-1.5 font-bold leading-none">{new Date(task.Date_Task || Date.now()).toLocaleDateString('id-ID', {month:'short', day:'numeric'})}</span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    );
  }

  // ── SISWA DASHBOARD ──────────────────────────────────────
  return (
    <div className="space-y-8 animate-fade-in-up">

      {/* Hero */}
      <div className="relative overflow-hidden rounded-3xl p-8 md:p-10"
        style={{
          background: 'linear-gradient(135deg, #0f2027 0%, #1a3a2a 50%, #0d2318 100%)',
          border: '1px solid rgba(52,211,153,0.2)',
          boxShadow: '0 20px 60px rgba(0,0,0,0.4), inset 0 1px 0 rgba(255,255,255,0.04)',
        }}>
        <div className="absolute top-0 right-0 w-80 h-80 rounded-full pointer-events-none"
          style={{ background: 'radial-gradient(circle, rgba(52,211,153,0.2), transparent 65%)', transform: 'translate(30%,-30%)' }} />
        <div className="absolute inset-0 opacity-5 pointer-events-none"
          style={{ backgroundImage: 'radial-gradient(circle, white 1px, transparent 1px)', backgroundSize: '32px 32px' }} />
        <div className="relative z-10 flex items-center gap-6 md:gap-8">
          <div className="w-16 h-16 md:w-20 md:h-20 rounded-2xl flex items-center justify-center text-2xl md:text-3xl font-black text-white shrink-0 shadow-2xl"
            style={{ background: 'linear-gradient(135deg, rgba(52,211,153,0.5), rgba(16,185,129,0.5))', border: '2px solid rgba(255,255,255,0.15)' }}>
            {user?.username?.charAt(0).toUpperCase() || 'S'}
          </div>
          <div className="space-y-1">
            <p className="text-xs md:text-sm font-bold tracking-wide leading-none" style={{ color: '#6ee7b7' }}>{greeting.emoji} {greeting.text},</p>
            <h1 className="text-2xl md:text-4xl font-black text-white leading-tight">{user?.username || 'Siswa'}</h1>
            <div className="flex items-center gap-3 pt-2 flex-wrap">
              <span className="badge badge-success px-4 py-1 font-bold">👨‍🎓 Siswa</span>
              {studentData && <span className="badge badge-info px-4 py-1 font-bold">📚 {studentData.kelas}</span>}
              <span className="badge px-4 py-1 font-bold" style={{ background: 'rgba(255,255,255,0.08)', color: '#94a3b8', border: '1px solid rgba(255,255,255,0.12)' }}>
                📅 {today}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Student Profile Card */}
      {loadingStats ? (
        <div className="glass rounded-3xl p-8">
          <div className="grid grid-cols-2 md:grid-cols-3 gap-6 md:gap-8">
            {[1,2,3,4,5,6].map(i => (
              <div key={i} className="h-16 md:h-20 rounded-2xl bg-white/5 animate-pulse" />
            ))}
          </div>
        </div>
      ) : studentData ? (
        <div className="glass rounded-3xl overflow-hidden shadow-2xl border border-white/5">
          <div className="px-6 py-5 md:px-8 md:py-6" style={{ borderBottom: '1px solid rgba(255,255,255,0.06)', background: 'rgba(255,255,255,0.01)' }}>
            <h3 className="font-black text-white text-lg flex items-center gap-3 leading-tight">
              <GraduationCap size={20} style={{ color: '#34d399' }} /> Profil Siswa
            </h3>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 divide-y sm:divide-y-0 sm:divide-x divide-white/5">
            {[
              { label: 'Nama Lengkap', value: studentData.full_name, icon: '👤' },
              { label: 'Kelas', value: studentData.kelas, icon: '🏫' },
              { label: 'Jurusan', value: studentData.jurusan, icon: '📖' },
              { label: 'No. Absen', value: studentData.absen.toString(), icon: '🔢' },
              { label: 'Wali Kelas', value: studentData.wali_kelas, icon: '👨‍🏫' },
              { label: 'Mata Pelajaran', value: studentData.mapel_students, icon: '📚' },
            ].map((item, i) => (
              <div key={i} className="p-6 md:p-8 space-y-3 bg-[#0d1117]/50 hover:bg-[#111827] transition-all group">
                <p className="text-[10px] md:text-xs font-black uppercase tracking-widest leading-none text-slate-500 group-hover:text-emerald-400 transition-colors">{item.icon} {item.label}</p>
                <p className="text-sm md:text-base font-black text-white truncate leading-snug">{item.value || '—'}</p>
              </div>
            ))}
          </div>
        </div>
      ) : (
        <div className="glass rounded-3xl p-12 md:p-20 text-center border-2 border-dashed border-white/10 space-y-6">
          <div className="w-16 h-16 md:w-20 md:h-20 rounded-3xl flex items-center justify-center mx-auto shadow-xl"
            style={{ background: 'rgba(52,211,153,0.1)', border: '1px solid rgba(52,211,153,0.2)' }}>
            <GraduationCap size={32} className="md:w-8 md:h-8 text-emerald-400" />
          </div>
          <div className="space-y-2">
            <h3 className="font-black text-white text-xl leading-tight">Belum Ada Profil Siswa</h3>
            <p className="text-xs md:text-sm text-slate-500 leading-relaxed max-w-xs mx-auto">Silakan lengkapi profil akademik Anda untuk mulai menggunakan fitur dashboard secara penuh.</p>
          </div>
          <Link href="/dashboard/students">
            <button className="btn-primary shadow-lg shadow-emerald-500/20">
              Lengkapi Sekarang →
            </button>
          </Link>
        </div>
      )}

      {/* Quick links */}
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-6 md:gap-8">
        {[
          { label: 'Tugas Saya', href: '/dashboard/tasks', icon: BookOpen, color: '#f472b6', glow: 'rgba(244,114,182,0.15)' },
          { label: 'Absensi', href: '/dashboard/absensi', icon: CalendarCheck, color: '#fbbf24', glow: 'rgba(251,191,36,0.15)' },
          { label: 'Nilai', href: '/dashboard/grades', icon: Award, color: '#818cf8', glow: 'rgba(129,140,248,0.15)' },
          { label: 'Profil Saya', href: '/dashboard/students', icon: Star, color: '#34d399', glow: 'rgba(52,211,153,0.15)' },
        ].map((q, i) => (
          <Link key={i} href={q.href}>
            <div className="glass glass-hover rounded-2xl p-6 md:p-8 flex flex-col items-center gap-4 cursor-pointer text-center group">
              <div className="w-12 h-12 md:w-14 md:h-14 rounded-2xl flex items-center justify-center shadow-lg transition-transform group-hover:scale-110"
                style={{ background: q.glow }}>
                <q.icon size={24} style={{ color: q.color }} />
              </div>
              <span className="text-[11px] md:text-xs font-black uppercase tracking-widest text-white leading-none">{q.label}</span>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}
