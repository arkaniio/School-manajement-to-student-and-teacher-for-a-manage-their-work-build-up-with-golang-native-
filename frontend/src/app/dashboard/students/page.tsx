"use client";

import React, { useEffect, useState } from 'react';
import api from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Plus, Edit, Trash2, GraduationCap, CheckCircle, UserCheck, Users, BookOpen, User, MapPin, School, BadgeCheck } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';

// ── Types ────────────────────────────────────────────────────────────────────

interface StudentData {
  id: string;
  full_name: string;
  kelas: string;
  jurusan: string;
  absen: number;
  wali_kelas: string;
  mapel_students: string;
  student_profile: string;
  user_id: string;
}

// ── Main Component ────────────────────────────────────────────────────────────

export default function StudentsPage() {
  const { user } = useAuth();
  const isSiswa = user?.role === 'siswa';
  const isGuru = user?.role === 'guru';

  const [student, setStudent] = useState<StudentData | null>(null);
  const [allStudents, setAllStudents] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    full_name: '', kelas: '', jurusan: '', absen: 0,
    student_profile: 'default.jpg', wali_kelas: '', mapel_students: ''
  });
  const [formLoading, setFormLoading] = useState(false);
  const [formError, setFormError] = useState('');

  const fetchStudentInfo = async () => {
    try {
      setLoading(true);
      if (isSiswa) {
        const res = await api.get('/students');
        setStudent(res.data?.data || null);
      } else if (isGuru) {
        // As specified in requirements, use tasks data to list students for guru
        const res = await api.get('/tasks');
        const tasks = res.data?.data || [];
        
        const studentMap = new Map();
        tasks.forEach((t: any) => {
          if (t.Student_Id && t.Students && (t.Students.Full_Name || t.Students.full_name)) {
             if (!studentMap.has(t.Student_Id)) {
               studentMap.set(t.Student_Id, {
                 id: t.Student_Id,
                 name: t.Students.Full_Name || t.Students.full_name,
                 kelas: t.Students.Kelas || t.Students.kelas || '—',
                 jurusan: t.Students.Jurusan || t.Students.jurusan || '—'
               });
             }
          }
        });
        setAllStudents(Array.from(studentMap.values()));
      }
    } catch (err: any) {
      if (err?.code !== 'ERR_CANCELED') console.error("Failed to fetch data", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStudentInfo();
  }, [isSiswa, isGuru]);

  const openAddModal = () => {
    setEditingId(null);
    setFormError('');
    setFormData({ 
      full_name: user?.username || '', 
      kelas: '', 
      jurusan: '', 
      absen: 0, 
      student_profile: 'default.jpg', 
      wali_kelas: '', 
      mapel_students: '' 
    });
    setIsModalOpen(true);
  };

  const openEditModal = () => {
    if (!student) return;
    setEditingId(student.id);
    setFormError('');
    setFormData({
      full_name: student.full_name || '',
      kelas: student.kelas || '',
      jurusan: student.jurusan || '',
      absen: student.absen || 0,
      student_profile: student.student_profile || '',
      wali_kelas: student.wali_kelas || '',
      mapel_students: student.mapel_students || '',
    });
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError('');
    setFormLoading(true);
    try {
      const payload = { ...formData, absen: parseInt(formData.absen.toString(), 10) };
      if (editingId) {
        await api.patch(`/student/${editingId}`, payload);
      } else {
        await api.post('/student', payload);
      }
      setIsModalOpen(false);
      fetchStudentInfo();
    } catch (err: any) {
      setFormError(err.response?.data?.message || 'Failed to save student data.');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!student || !confirm('Are you sure you want to delete your student profile?')) return;
    try {
      await api.delete(`/student/${student.id}`);
      setStudent(null);
    } catch (err) {
      alert('Failed to delete student profile.');
    }
  };

  // ── GURU VIEW ─────────────────────────────────────────────────────────────
  if (isGuru) {
    return (
      <div className="space-y-8 animate-fade-in-up">
        <div className="space-y-2">
          <h1 className="text-2xl font-black text-white flex items-center gap-3 uppercase tracking-tight">
            <Users size={24} className="text-indigo-400" /> Direktori Siswa
          </h1>
          <p className="text-slate-400 text-sm font-medium leading-relaxed">Lihat daftar siswa yang telah mendaftar dan mengumpulkan tugas.</p>
        </div>
        
        <div className="glass rounded-3xl border border-white/5 shadow-2xl overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-white/1 border-b border-white/5">
                  <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em]">Student ID</th>
                  <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em]">Nama Lengkap</th>
                  <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em]">Kelas</th>
                  <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em]">Jurusan</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/5">
                {loading ? (
                  <tr>
                    <td colSpan={4} className="px-8 py-20 text-center">
                      <div className="animate-spin w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full mx-auto mb-4" />
                      <p className="text-xs text-slate-500 font-bold uppercase tracking-widest">Loading students...</p>
                    </td>
                  </tr>
                ) : allStudents.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="px-8 py-20 text-center space-y-4">
                      <div className="w-16 h-16 bg-white/5 rounded-2xl flex items-center justify-center mx-auto text-slate-600">
                        <GraduationCap size={32} />
                      </div>
                      <p className="font-black text-lg text-slate-400 uppercase tracking-tight">Belum ada data siswa.</p>
                    </td>
                  </tr>
                ) : (
                  allStudents.map((s, i) => (
                    <tr key={i} className="hover:bg-white/2 transition-all group">
                      <td className="px-8 py-6">
                        <span className="text-xs font-black text-slate-500 font-mono tracking-tight uppercase px-3 py-1 bg-white/5 rounded-lg border border-white/5">
                          #{s.id.split('-')[0].toUpperCase()}
                        </span>
                      </td>
                      <td className="px-8 py-6">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-xl bg-linear-to-br from-indigo-500/20 to-purple-500/20 border border-white/10 flex items-center justify-center text-indigo-300 font-black text-sm">
                            {s.name?.charAt(0) || <User size={16} />}
                          </div>
                          <span className="text-sm font-black text-white group-hover:text-indigo-400 transition-colors uppercase">
                            {s.name}
                          </span>
                        </div>
                      </td>
                      <td className="px-8 py-6">
                        <span className="px-3 py-1 bg-blue-500/10 text-blue-400 border border-blue-500/20 rounded-full text-[10px] font-black uppercase tracking-widest leading-none">
                          {s.kelas}
                        </span>
                      </td>
                      <td className="px-8 py-6">
                        <span className="px-3 py-1 bg-purple-500/10 text-purple-400 border border-purple-500/20 rounded-full text-[10px] font-black uppercase tracking-widest leading-none">
                          {s.jurusan}
                        </span>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    );
  }

  // ── SISWA VIEW ────────────────────────────────────────────────────────────
  return (
    <div className="space-y-8 animate-fade-in-up">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
        <div className="space-y-2">
          <h1 className="text-2xl font-black text-white uppercase tracking-tight">Profil Siswa</h1>
          <p className="text-slate-400 text-sm font-medium leading-relaxed">Kelola data registrasi dan informasi akademik kamu.</p>
        </div>
        {isSiswa && !loading && !student && (
          <Button 
            onClick={openAddModal}
            icon={<Plus size={18} />}
            className="font-black tracking-wide"
          >
            Register as Student
          </Button>
        )}
      </div>

      {loading ? (
        <div className="glass rounded-3xl p-16 text-center space-y-4">
          <div className="animate-spin w-10 h-10 border-4 border-indigo-500 border-t-transparent rounded-full mx-auto" />
          <p className="text-slate-400 text-sm font-bold uppercase tracking-widest">Loading student profile...</p>
        </div>
      ) : student ? (
        <div className="space-y-8">
          <div className="flex items-center gap-5 p-5 rounded-3xl bg-emerald-500/5 border border-emerald-500/20 text-emerald-400 animate-fade-in overflow-hidden relative">
            <div className="absolute top-0 right-0 w-32 h-32 bg-emerald-500/10 blur-3xl rounded-full translate-x-1/2 -translate-y-1/2 pointer-events-none" />
            <div className="w-12 h-12 rounded-2xl bg-emerald-500/10 flex items-center justify-center shrink-0 border border-emerald-500/10">
              <CheckCircle size={24} />
            </div>
            <div className="space-y-1 relative">
              <p className="font-black text-sm uppercase tracking-wider leading-tight">Profile Terverifikasi</p>
              <p className="text-xs font-medium opacity-80 leading-relaxed">Kamu telah terdaftar sebagai siswa resmi. Profil kamu dapat diakses oleh guru.</p>
            </div>
          </div>

          <div className="glass rounded-3xl overflow-hidden shadow-2xl border border-white/5">
            <div className="relative p-8 md:p-12 flex flex-col md:flex-row items-center gap-6 md:gap-10 overflow-hidden"
              style={{ background: 'linear-gradient(135deg, #1e1b4b 0%, #312e81 100%)' }}>
              
              {/* Background Decoration */}
              <div className="absolute inset-0 opacity-10"
                style={{ backgroundImage: 'radial-gradient(circle, white 1px, transparent 1px)', backgroundSize: '24px 24px' }} />
              <div className="absolute -bottom-20 -left-20 w-64 h-64 bg-indigo-500/20 rounded-full blur-[100px] pointer-events-none" />
              
              <div className="relative w-24 h-24 md:w-32 md:h-32 rounded-3xl bg-white/10 border-2 border-white/20 flex items-center justify-center text-white text-4xl font-black shrink-0 shadow-2xl group transition-transform hover:scale-105 duration-500 border-dashed">
                {student.full_name?.charAt(0)?.toUpperCase() || <User size={48} />}
              </div>
              
              <div className="relative text-white text-center md:text-left space-y-4 flex-1">
                <div className="space-y-2">
                  <div className="flex items-center justify-center md:justify-start gap-3">
                    <h2 className="text-3xl md:text-4xl font-black tracking-tighter uppercase leading-tight">{student.full_name}</h2>
                    <BadgeCheck size={20} className="text-indigo-400 shrink-0" />
                  </div>
                  <p className="text-indigo-200/60 text-xs font-black uppercase tracking-[0.2em] leading-none flex items-center justify-center md:justify-start gap-2">
                    <span className="w-1.5 h-1.5 rounded-full bg-indigo-500" /> ID: #{student.id.split('-')[0].toUpperCase()}
                  </p>
                </div>
                <div className="flex flex-wrap justify-center md:justify-start gap-3">
                  <span className="bg-white/5 backdrop-blur-xl border border-white/10 text-[10px] font-black uppercase tracking-widest px-5 py-2 rounded-full leading-none flex items-center gap-2">
                    <School size={12} className="text-indigo-400" /> Kelas {student.kelas}
                  </span>
                  <span className="bg-white/5 backdrop-blur-xl border border-white/10 text-[10px] font-black uppercase tracking-widest px-5 py-2 rounded-full leading-none flex items-center gap-2">
                    <GraduationCap size={12} className="text-purple-400" /> {student.jurusan}
                  </span>
                </div>
              </div>

              <div className="relative flex gap-3 shrink-0 pt-4 md:pt-0">
                <button onClick={openEditModal}
                  className="p-4 bg-white/5 hover:bg-white/10 rounded-2xl border border-white/10 text-white transition-all shadow-lg active:scale-90 group" title="Edit Profile">
                  <Edit size={20} className="group-hover:text-indigo-400 transition-colors" />
                </button>
                <button onClick={handleDelete}
                  className="p-4 bg-red-500/10 hover:bg-red-500/20 rounded-2xl border border-red-500/20 text-red-500 transition-all shadow-lg active:scale-90 group" title="Delete Profile">
                  <Trash2 size={20} className="group-hover:scale-110 transition-transform" />
                </button>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 divide-y sm:divide-y-0 sm:divide-x divide-white/5 bg-white/1">
              {[
                { label: 'Nomor Absen', value: student.absen?.toString(), icon: <span className="text-indigo-400 font-black">#</span>, badge: 'badge-info' },
                { label: 'Wali Kelas', value: student.wali_kelas, icon: <User size={16} className="text-emerald-400" />, badge: 'badge-success' },
                { label: 'Mata Pelajaran', value: student.mapel_students, icon: <BookOpen size={16} className="text-purple-400" />, badge: 'badge-purple' },
              ].map((item, i) => (
                <div key={i} className="p-8 space-y-4 hover:bg-white/2 transition-colors group">
                  <div className="flex items-center gap-3">
                    {item.icon}
                    <p className="text-[10px] font-black text-slate-500 uppercase tracking-[0.2em] leading-none">{item.label}</p>
                  </div>
                  <p className="text-lg font-black text-white leading-tight uppercase group-hover:translate-x-1 transition-transform">{item.value || '—'}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      ) : (
        <div className="glass rounded-3xl p-12 md:p-24 text-center border-dashed border-2 border-white/5 space-y-8 bg-white/1">
          <div className="w-24 h-24 mx-auto rounded-3xl flex items-center justify-center bg-indigo-500/5 border border-indigo-500/10 shadow-inner group-hover:scale-110 transition-transform duration-500">
            <UserCheck size={48} className="text-indigo-500/40" />
          </div>
          <div className="space-y-3">
            <h3 className="text-2xl font-black text-white leading-tight uppercase tracking-tight">Profil Belum Terdaftar</h3>
            <p className="text-slate-500 text-sm font-medium leading-relaxed max-w-sm mx-auto">Silakan daftarkan profil siswa kamu untuk dapat melihat dashboard, mengumpulkan tugas, dan melihat nilai akademik.</p>
          </div>
          <Button 
            onClick={openAddModal} 
            variant="primary" 
            size="lg" 
            icon={<Plus size={20} />}
            className="font-black tracking-widest uppercase"
          >
            Daftar Sekarang
          </Button>
        </div>
      )}

      {/* Modal Profile */}
      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        title={editingId ? "Edit Profil Siswa" : "Registrasi Siswa Baru"}
      >
        <form onSubmit={handleSubmit} className="space-y-6">
          {formError && (
            <div className="p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-2xl text-sm font-semibold leading-normal flex items-center gap-3 animate-head-shake">
              <span className="text-lg">⚠</span> {formError}
            </div>
          )}
          
          <div className="space-y-4">
            <Input 
              label="Nama Lengkap" 
              value={formData.full_name}
              onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
              placeholder="Masukkan nama lengkap kamu" 
              required 
              icon={<User size={18} />}
            />
            
            <div className="grid grid-cols-2 gap-6">
              <Input 
                label="Kelas" 
                value={formData.kelas}
                onChange={(e) => setFormData({ ...formData, kelas: e.target.value })}
                placeholder="Contoh: XII" 
                required 
                icon={<School size={18} />}
              />
              <Input 
                label="Jurusan" 
                value={formData.jurusan}
                onChange={(e) => setFormData({ ...formData, jurusan: e.target.value })}
                placeholder="Contoh: IPA" 
                required 
                icon={<GraduationCap size={18} />}
              />
            </div>
            
            <div className="grid grid-cols-2 gap-6">
              <Input 
                label="Nomor Absen" 
                type="number" 
                value={formData.absen}
                onChange={(e) => setFormData({ ...formData, absen: parseInt(e.target.value) || 0 })}
                required 
                placeholder="0"
              />
              <Input 
                label="Wali Kelas" 
                value={formData.wali_kelas}
                onChange={(e) => setFormData({ ...formData, wali_kelas: e.target.value })}
                placeholder="Nama Bapak/Ibu Guru" 
                required 
                icon={<UserCheck size={18} />}
              />
            </div>
            
            <Input 
              label="Mata Pelajaran Favorit" 
              value={formData.mapel_students}
              onChange={(e) => setFormData({ ...formData, mapel_students: e.target.value })}
              placeholder="e.g. Matematika, Fisika, Biologi" 
              required 
              icon={<BookOpen size={18} />}
            />
          </div>

          <div className="pt-6 flex justify-end gap-4 border-t border-white/5">
            <Button 
              type="button" 
              variant="secondary"
              onClick={() => setIsModalOpen(false)}
              className="font-bold tracking-tight"
            >
              Cancel
            </Button>
            <Button 
              type="submit" 
              isLoading={formLoading}
              className="font-black tracking-wide"
            >
              {editingId ? 'Simpan Perubahan' : 'Selesaikan Registrasi'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
