"use client";

import React, { useEffect, useState } from 'react';
import api from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { Calendar, UserCheck, UserX, UserMinus, Plus, Trash2, Edit } from 'lucide-react';

export default function AttendancePage() {
  const { user } = useAuth();
  const isGuru = user?.role === 'guru';
  const isSiswa = user?.role === 'siswa';
  const [records, setRecords] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [studentId, setStudentId] = useState<string | null>(null);

  // Form State based on PayloadAbsensis
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({
    name_lengkap: '',
    kelas: '',
    jurusan: '',
    hari: '',
    tanggal: '',
    status: 'h',
    keterangan: '',
    keterangan_tidak_hadir: '',
    keterangan_dispen: '',
    file_dispen: ''
  });
  const [formLoading, setFormLoading] = useState(false);

  const fetchAttendance = async () => {
    try {
      setLoading(true);
      if (isGuru) {
        const res = await api.get('/absensi/all-with-students');
        setRecords(res.data?.data || []);
      } else if (isSiswa) {
        // also fetch student_id needed for creating absensi
        const [absensiRes, studentRes] = await Promise.all([
          api.get('/absensi/all-with-students'),
          api.get('/students'),
        ]);
        
        const sid = studentRes.data?.data?.id || studentRes.data?.data?.Id;
        if (sid) {
          setStudentId(sid);
          const allRecords = absensiRes.data?.data || [];
          setRecords(allRecords.filter((r: any) => r.StudentId === sid || r.student_id === sid || r.Student_Id === sid));
        } else {
          setRecords([]);
        }
      }
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAttendance();
  }, [isGuru, isSiswa]);

  const openAddModal = () => {
    setEditingId(null);
    setFormData({ 
      name_lengkap: '',
      kelas: '',
      jurusan: '',
      hari: 'Senin',
      tanggal: new Date().toISOString().split('T')[0],
      status: 'accepted',
      keterangan: 'hadir',
      keterangan_tidak_hadir: '',
      keterangan_dispen: '',
      file_dispen: ''
    });
    setIsModalOpen(true);
  };

  const openEditModal = (record: any) => {
    setEditingId(record.Id || record.id);
    setFormData({ 
      name_lengkap: record.NameLengkap || record.name_lengkap || record.StudentFullName || '', 
      kelas: record.KelasAbsensi || record.kelas_absensi || '', 
      jurusan: record.JurusanAbsensi || record.jurusan_absensi || '',
      hari: record.Hari || record.hari || '',
      tanggal: record.Tanggal ? new Date(record.Tanggal).toISOString().split('T')[0] : '', 
      status: record.Status || record.status || 'accepted', 
      keterangan: record.Keterangan || record.keterangan || 'hadir',
      keterangan_tidak_hadir: record.KeteranganTidakHadir || record.keterangan_tidak_hadir || '',
      keterangan_dispen: record.KeteranganDispen || record.keterangan_dispen || '',
      file_dispen: record.FileDispen || record.file_dispen || ''
    });
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormLoading(true);
    try {
      const payload = {
        ...formData,
        student_id: studentId || '',
        tanggal: new Date(formData.tanggal).toISOString(),
      };

      if (editingId) {
        await api.patch(`/absensi/${editingId}`, payload);
      } else {
        if (!studentId) {
          alert('Data siswa belum tersedia. Pastikan kamu sudah terdaftar sebagai siswa.');
          return;
        }
        await api.post('/absensi', payload);
      }
      setIsModalOpen(false);
      fetchAttendance();
    } catch (err: any) {
      console.error(err);
      alert(err.response?.data?.message || 'Failed to save attendance.');
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this record?')) return;
    try {
      await api.delete(`/absensi/${id}`);
      setRecords(records.filter(r => (r.Id || r.id) !== id));
    } catch (err) {
      console.error(err);
      alert('Failed to delete attendance record.');
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status?.toLowerCase()) {
      case 'accepted':
      case 'h':
        return <span className="flex w-fit items-center gap-2 px-3 py-1.5 rounded-full text-[11px] font-black uppercase tracking-wider bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"><UserCheck size={14}/> Hadir</span>;
      case 'permissions':
      case 'i':
        return <span className="flex w-fit items-center gap-2 px-3 py-1.5 rounded-full text-[11px] font-black uppercase tracking-wider bg-blue-500/10 text-blue-400 border border-blue-500/20"><Calendar size={14}/> Izin</span>;
      case 'not accepted':
      case 's':
      case 'a':
        return <span className="flex w-fit items-center gap-2 px-3 py-1.5 rounded-full text-[11px] font-black uppercase tracking-wider bg-red-500/10 text-red-400 border border-red-500/20"><UserX size={14}/> Absen</span>;
      default:
        return <span className="flex w-fit px-3 py-1.5 rounded-full text-[11px] font-black uppercase tracking-wider bg-white/5 text-slate-400">? ({status})</span>;
    }
  };

  return (
    <div className="space-y-8 animate-fade-in-up">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
        <div className="space-y-2">
          <h1 className="text-2xl font-bold text-white leading-tight flex items-center gap-3">
            <Calendar size={24} className="text-blue-400" /> Attendance Records
          </h1>
          <p className="text-slate-400 text-sm leading-relaxed">Kelola dan lihat riwayat absensi harian kelas.</p>
        </div>
        
        {isSiswa && (
          <button onClick={openAddModal} className="btn-primary shrink-0 transition-transform active:scale-95 shadow-lg shadow-indigo-500/20">
            <Plus size={18} /> Record Attendance
          </button>
        )}
      </div>

      <div className="glass rounded-3xl border border-white/5 shadow-2xl overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm text-slate-300">
            <thead className="bg-white/[.02] text-slate-500 font-black uppercase tracking-widest text-[10px] border-b border-white/5">
              <tr>
                <th className="px-8 py-5">Date & Day</th>
                <th className="px-8 py-5">Student Name</th>
                <th className="px-8 py-5">Status</th>
                <th className="px-8 py-5 hidden md:table-cell">Keterangan</th>
                <th className="px-8 py-5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {loading ? (
                <tr>
                  <td colSpan={5} className="px-8 py-16 text-center text-slate-500 leading-relaxed italic">
                    Loading attendance records...
                  </td>
                </tr>
              ) : records.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-8 py-20 text-center text-slate-500 leading-relaxed">
                    <Calendar size={48} className="mx-auto mb-4 opacity-20" />
                    <p className="font-bold text-lg text-slate-400">No attendance records found.</p>
                    <p className="text-sm">Silakan lakukan absensi terlebih dahulu.</p>
                  </td>
                </tr>
              ) : (
                records.map((record, i) => {
                  const id = record.Id || record.id || i;
                  return (
                  <tr key={id} className="hover:bg-white/[.02] group transition-all">
                    <td className="px-8 py-6 font-bold text-white leading-relaxed">
                      <div className="flex flex-col">
                        <span>{record.Tanggal ? new Date(record.Tanggal).toLocaleDateString() : 'No date'}</span>
                        <span className="text-[10px] text-slate-500 uppercase tracking-widest">{record.Hari || record.hari}</span>
                      </div>
                    </td>
                    <td className="px-8 py-6 font-black text-slate-100 leading-relaxed">
                      {record.NameLengkap || record.name_lengkap || record.StudentFullName || 'Unknown Student'}
                    </td>
                    <td className="px-8 py-6">
                      {getStatusBadge(record.Status || record.status)}
                    </td>
                    <td className="px-8 py-6 hidden md:table-cell text-slate-400 text-xs leading-relaxed max-w-xs truncate">
                      {record.Keterangan || record.keterangan || '—'}
                    </td>
                    <td className="px-8 py-6 text-right">
                      <div className="flex justify-end gap-3 opacity-0 group-hover:opacity-100 transition-all transform translate-x-2 group-hover:translate-x-0">
                        <button 
                          className="p-2 text-slate-500 hover:text-white hover:bg-white/10 rounded-xl transition-all" title="Edit"
                          onClick={() => openEditModal(record)}
                        >
                          <Edit size={16} />
                        </button>
                        <button 
                          className="p-2 text-slate-500 hover:text-red-400 hover:bg-red-400/10 rounded-xl transition-all" title="Delete"
                          onClick={() => handleDelete(id)}
                        >
                          <Trash2 size={16} />
                        </button>
                      </div>
                    </td>
                  </tr>
                )})
              )}
            </tbody>
          </table>
        </div>
      </div>

      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)} 
        title={editingId ? "Edit Attendance" : "Record Attendance"}
      >
        <form onSubmit={handleSubmit} className="space-y-6 max-h-[75vh] overflow-y-auto p-2 pr-4 custom-scrollbar">
          <div className="space-y-4">
            <Input 
              label="Nama Lengkap" 
              value={formData.name_lengkap} 
              onChange={e => setFormData({...formData, name_lengkap: e.target.value})} 
              required 
            />
            <div className="grid grid-cols-2 gap-6">
              <Input 
                label="Kelas" 
                value={formData.kelas} 
                onChange={e => setFormData({...formData, kelas: e.target.value})} 
                required 
              />
              <Input 
                label="Jurusan" 
                value={formData.jurusan} 
                onChange={e => setFormData({...formData, jurusan: e.target.value})} 
                required 
              />
            </div>
            <div className="grid grid-cols-2 gap-6">
              <Input 
                label="Hari" 
                value={formData.hari} 
                onChange={e => setFormData({...formData, hari: e.target.value})} 
                required 
              />
              <Input 
                label="Tanggal" 
                type="date"
                value={formData.tanggal} 
                onChange={e => setFormData({...formData, tanggal: e.target.value})} 
                required 
              />
            </div>
            <div className="space-y-2">
              <label className="block text-xs font-bold uppercase tracking-wider text-slate-400 leading-none">Keterangan Dasar</label>
              <select
                value={formData.keterangan}
                onChange={e => setFormData({...formData, keterangan: e.target.value})}
                className="input-dark leading-relaxed"
                required
              >
                <option value="hadir" className="bg-slate-900">Hadir</option>
                <option value="tidak hadir" className="bg-slate-900">Tidak Hadir (Sakit/Alpa)</option>
                <option value="izin" className="bg-slate-900">Izin</option>
                <option value="dispen" className="bg-slate-900">Dispensasi</option>
              </select>
            </div>
            {formData.keterangan !== 'hadir' && (
              <div className="space-y-4 pt-2 border-t border-white/5 mt-4">
                <Input 
                  label="Keterangan Tidak Hadir" 
                  value={formData.keterangan_tidak_hadir} 
                  onChange={e => setFormData({...formData, keterangan_tidak_hadir: e.target.value})} 
                  placeholder="e.g. Sakit perut / Keperluan keluarga"
                />
                <Input 
                  label="Keterangan Dispen" 
                  value={formData.keterangan_dispen} 
                  onChange={e => setFormData({...formData, keterangan_dispen: e.target.value})} 
                  placeholder="e.g. Lomba OSN"
                />
                <Input 
                  label="File Dispen (url/filename)" 
                  value={formData.file_dispen} 
                  onChange={e => setFormData({...formData, file_dispen: e.target.value})} 
                  placeholder="upload_file.jpg"
                />
              </div>
            )}
          </div>
          <div className="pt-6 flex justify-end gap-4 border-t border-white/5">
            <button type="button" onClick={() => setIsModalOpen(false)}
              className="px-6 py-3 text-sm font-bold rounded-xl bg-white/5 hover:bg-white/10 text-slate-300 transition-colors leading-none">
              Cancel
            </button>
            <button type="submit" disabled={formLoading}
              className="btn-primary shadow-lg shadow-indigo-500/20"
              style={{ opacity: formLoading ? 0.6 : 1 }}>
              {formLoading ? 'Saving...' : 'Save Record'}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
