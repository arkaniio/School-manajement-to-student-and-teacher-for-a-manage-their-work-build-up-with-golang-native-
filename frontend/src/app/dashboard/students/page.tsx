"use client";

import React, { useEffect, useState } from 'react';
import api from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Plus, Edit, Trash2, GraduationCap, CheckCircle, UserCheck, Users } from 'lucide-react';

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

const InputField = ({
  label, type = 'text', value, onChange, placeholder, required
}: {
  label: string; type?: string; value: any;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder?: string; required?: boolean;
}) => (
  <div className="space-y-2">
    <label className="block text-xs font-semibold uppercase tracking-wider text-slate-400 leading-none">{label}</label>
    <input
      type={type} value={value} onChange={onChange}
      placeholder={placeholder} required={required}
      className="input-dark"
    />
  </div>
);

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
        const res = await api.get('/tasks');
        const tasks = res.data?.data || [];
        
        const studentMap = new Map();
        tasks.forEach((t: any) => {
          if (t.Student_Id && t.Students && t.Students.Full_Name) {
             if (!studentMap.has(t.Student_Id)) {
               studentMap.set(t.Student_Id, {
                 id: t.Student_Id,
                 name: t.Students.Full_Name,
                 kelas: t.Students.Kelas || '—',
                 jurusan: t.Students.Jurusan || '—'
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
    setFormData({ full_name: '', kelas: '', jurusan: '', absen: 0, student_profile: 'default.jpg', wali_kelas: '', mapel_students: '' });
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
          <h1 className="text-2xl font-bold text-white flex items-center gap-3 leading-tight">
            <Users size={24} className="text-emerald-400" /> Student Directory
          </h1>
          <p className="text-slate-400 text-sm leading-relaxed">Lihat daftar siswa yang telah mendaftar dan mengumpulkan tugas.</p>
        </div>
        
        <div className="glass rounded-2xl border border-transparent shadow-[0_4px_20px_rgba(0,0,0,0.2)] overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-300">
              <thead className="bg-white/5 text-slate-300 font-medium border-b border-white/5">
                <tr>
                  <th className="px-6 py-4 leading-none">Siswa ID</th>
                  <th className="px-6 py-4 leading-none">Nama Lengkap</th>
                  <th className="px-6 py-4 leading-none">Kelas</th>
                  <th className="px-6 py-4 leading-none">Jurusan</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white/5">
                {loading ? (
                  <tr><td colSpan={4} className="px-6 py-10 text-center text-slate-400 leading-relaxed">Loading students...</td></tr>
                ) : allStudents.length === 0 ? (
                  <tr><td colSpan={4} className="px-6 py-16 text-center text-slate-400">
                    <div className="space-y-3">
                      <GraduationCap size={40} className="mx-auto text-slate-500" />
                      <p className="font-medium text-base">Belum ada data siswa.</p>
                    </div>
                  </td></tr>
                ) : (
                  allStudents.map((s, i) => (
                    <tr key={i} className="hover:bg-white/5 transition-colors group">
                      <td className="px-6 py-5 text-xs font-mono text-slate-500 leading-relaxed">{s.id.split('-')[0]}</td>
                      <td className="px-6 py-5 font-bold text-white leading-relaxed">{s.name}</td>
                      <td className="px-6 py-5 leading-relaxed"><span className="badge badge-info">{s.kelas}</span></td>
                      <td className="px-6 py-5 leading-relaxed"><span className="badge badge-purple">{s.jurusan}</span></td>
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
          <h1 className="text-2xl font-bold text-white leading-tight">My Student Profile</h1>
          <p className="text-slate-400 text-sm leading-relaxed">Manage your student registration data.</p>
        </div>
        {isSiswa && !loading && !student && (
          <button onClick={openAddModal}
            className="btn-primary shrink-0">
            <Plus size={16} /> Register as Student
          </button>
        )}
      </div>

      {loading ? (
        <div className="glass rounded-2xl p-16 text-center space-y-4">
          <div className="animate-spin w-10 h-10 border-4 border-indigo-500 border-t-transparent rounded-full mx-auto" />
          <p className="text-slate-400 text-sm leading-relaxed">Loading your student data...</p>
        </div>
      ) : student ? (
        <div className="space-y-8">
          <div className="flex items-center gap-5 p-5 rounded-2xl" style={{ background: 'rgba(16,185,129,0.1)', border: '1px solid rgba(16,185,129,0.2)', color: '#34d399' }}>
            <CheckCircle size={24} className="shrink-0" />
            <div className="space-y-1">
              <p className="font-bold text-sm leading-tight">Student Profile Registered</p>
              <p className="text-xs leading-relaxed" style={{ color: 'rgba(52,211,153,0.8)' }}>You are registered as a student. You can edit or delete your profile below.</p>
            </div>
          </div>

          <div className="glass rounded-2xl overflow-hidden shadow-xl">
            <div className="relative p-8 md:p-10 flex flex-col md:flex-row items-center gap-6 md:gap-8"
              style={{ background: 'linear-gradient(135deg, #312e81 0%, #4f46e5 100%)' }}>
              <div className="absolute inset-0 opacity-10"
                style={{ backgroundImage: 'radial-gradient(circle, white 1px, transparent 1px)', backgroundSize: '24px 24px' }} />
              
              <div className="relative w-20 h-20 md:w-24 md:h-24 rounded-3xl bg-white/20 border-2 border-white/30 flex items-center justify-center text-white text-3xl font-black shrink-0 shadow-2xl overflow-hidden">
                {student.full_name?.charAt(0)?.toUpperCase() || 'S'}
              </div>
              
              <div className="relative text-white text-center md:text-left space-y-4 flex-1">
                <div className="space-y-1">
                  <h2 className="text-2xl md:text-3xl font-black truncate leading-tight">{student.full_name}</h2>
                  <p className="text-indigo-100/80 text-sm font-medium leading-relaxed">Student ID: {student.id.split('-')[0]}</p>
                </div>
                <div className="flex flex-wrap justify-center md:justify-start gap-3">
                  <span className="bg-white/10 backdrop-blur-md border border-white/20 text-xs font-bold px-4 py-1.5 rounded-full leading-none">{student.kelas}</span>
                  <span className="bg-white/10 backdrop-blur-md border border-white/20 text-xs font-bold px-4 py-1.5 rounded-full leading-none">{student.jurusan}</span>
                </div>
              </div>

              <div className="relative flex gap-3 shrink-0">
                <button onClick={openEditModal}
                  className="p-3 bg-black/20 hover:bg-black/40 rounded-2xl border border-white/20 text-white transition-all hover:scale-105 active:scale-95" title="Edit">
                  <Edit size={20} />
                </button>
                <button onClick={handleDelete}
                  className="p-3 bg-red-500/30 hover:bg-red-500/50 rounded-2xl border border-white/20 text-white transition-all hover:scale-105 active:scale-95" title="Delete">
                  <Trash2 size={20} />
                </button>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 divide-y sm:divide-y-0 sm:divide-x divide-white/5 bg-white/[.01]">
              {[
                { label: 'No. Absen', value: student.absen?.toString(), emoji: '🔢' },
                { label: 'Wali Kelas', value: student.wali_kelas, emoji: '👨‍🏫' },
                { label: 'Mata Pelajaran', value: student.mapel_students, emoji: '📚' },
              ].map((item, i) => (
                <div key={i} className="p-6 md:p-8 space-y-3 hover:bg-white/[.02] transition-colors">
                  <p className="text-xs text-slate-400 font-bold tracking-widest uppercase leading-none flex items-center gap-2">
                    <span className="text-lg">{item.emoji}</span> {item.label}
                  </p>
                  <p className="text-lg font-black text-white leading-snug truncate">{item.value || '—'}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      ) : (
        <div className="glass rounded-3xl p-12 md:p-20 text-center border-dashed border-2 border-white/10 space-y-6">
          <div className="w-20 h-20 mx-auto rounded-3xl flex items-center justify-center bg-indigo-500/10 border border-indigo-500/20">
            <UserCheck size={40} className="text-indigo-400" />
          </div>
          <div className="space-y-2">
            <h3 className="text-xl font-black text-white leading-tight">No Student Profile Found</h3>
            <p className="text-slate-400 text-sm leading-relaxed max-w-xs mx-auto">You haven&apos;t registered your student data yet. Register now to access all features.</p>
          </div>
          <button onClick={openAddModal} className="btn-primary">
            <Plus size={16} /> Register as Student
          </button>
        </div>
      )}

      {/* Modal */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 md:p-6">
          <div className="absolute inset-0 bg-black/80 backdrop-blur-md" onClick={() => setIsModalOpen(false)} />
          <div className="relative glass rounded-3xl shadow-2xl w-full max-w-lg overflow-hidden border border-white/10" onClick={e => e.stopPropagation()}>
            <div className="h-1.5 rounded-t-3xl" style={{ background: 'linear-gradient(90deg, #10b981, #3b82f6)' }} />
            <div className="flex items-center justify-between px-8 py-6 border-b border-white/5 bg-white/[.02]">
              <div className="space-y-1">
                <h3 className="text-xl font-bold text-white leading-tight">{editingId ? 'Edit Student Profile' : 'Register as Student'}</h3>
                <p className="text-xs text-slate-400 leading-none">Complete the form below to save your data.</p>
              </div>
              <button onClick={() => setIsModalOpen(false)} className="text-slate-400 hover:text-white hover:bg-white/10 p-2.5 rounded-xl transition-colors">✕</button>
            </div>
            <form onSubmit={handleSubmit} className="p-8 space-y-6 overflow-y-auto max-h-[75vh]">
              {formError && (
                <div className="p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-xl text-sm leading-normal flex items-center gap-3">
                  <span className="text-lg">⚠</span> {formError}
                </div>
              )}
              
              <InputField label="Full Name" value={formData.full_name}
                onChange={(e) => setFormData({ ...formData, full_name: e.target.value })}
                placeholder="Enter your full name" required />
              
              <div className="grid grid-cols-2 gap-6">
                <InputField label="Kelas" value={formData.kelas}
                  onChange={(e) => setFormData({ ...formData, kelas: e.target.value })}
                  placeholder="e.g. XII" required />
                <InputField label="Jurusan" value={formData.jurusan}
                  onChange={(e) => setFormData({ ...formData, jurusan: e.target.value })}
                  placeholder="e.g. IPA" required />
              </div>
              
              <div className="grid grid-cols-2 gap-6">
                <InputField label="No. Absen" type="number" value={formData.absen}
                  onChange={(e) => setFormData({ ...formData, absen: parseInt(e.target.value) || 0 })}
                  required />
                <InputField label="Wali Kelas" value={formData.wali_kelas}
                  onChange={(e) => setFormData({ ...formData, wali_kelas: e.target.value })}
                  placeholder="Name of class teacher" required />
              </div>
              
              <InputField label="Mata Pelajaran" value={formData.mapel_students}
                onChange={(e) => setFormData({ ...formData, mapel_students: e.target.value })}
                placeholder="e.g. Matematika, Fisika" required />

              <div className="pt-6 flex justify-end gap-4 border-t border-white/5">
                <button type="button" onClick={() => setIsModalOpen(false)}
                  className="px-6 py-3 text-sm font-bold rounded-xl bg-white/5 hover:bg-white/10 text-slate-300 transition-colors leading-none">
                  Cancel
                </button>
                <button type="submit" disabled={formLoading}
                  className="btn-primary shadow-lg shadow-emerald-500/20" style={{ opacity: formLoading ? 0.6 : 1 }}>
                  {formLoading ? 'Saving...' : 'Save Profile'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
