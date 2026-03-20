"use client";

import React, { useEffect, useState } from 'react';
import api from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Modal } from '@/components/ui/Modal';
import { Calendar, UserCheck, UserX, UserMinus, Plus, Trash2, Edit, Clock, User, CheckCircle, Filter, Search } from 'lucide-react';

// ── Types ────────────────────────────────────────────────────────────────────

interface AbsensiRecord {
  Id?: string;
  id?: string;
  NameLengkap?: string;
  name_lengkap?: string;
  StudentFullName?: string;
  KelasAbsensi?: string;
  kelas_absensi?: string;
  JurusanAbsensi?: string;
  jurusan_absensi?: string;
  Hari?: string;
  hari?: string;
  Tanggal?: string;
  tanggal?: string;
  Status?: string;
  status?: string;
  Keterangan?: string;
  keterangan?: string;
  KeteranganTidakHadir?: string;
  keterangan_tidak_hadir?: string;
  KeteranganDispen?: string;
  keterangan_dispen?: string;
  FileDispen?: string;
  file_dispen?: string;
  StudentId?: string;
}

// ── Main Component ────────────────────────────────────────────────────────────

export default function AttendancePage() {
  const { user } = useAuth();
  const isGuru = user?.role === 'guru';
  const isSiswa = user?.role === 'siswa';
  const [records, setRecords] = useState<AbsensiRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [studentId, setStudentId] = useState<string | null>(null);

  // Form State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingId, setEditingId] = useState<string | null>(null);
  const [formData, setFormData] = useState({
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
  const [formLoading, setFormLoading] = useState(false);
  const [error, setError] = useState('');

  const fetchAttendance = async () => {
    try {
      setLoading(true);
      if (isGuru) {
        const res = await api.get('/absensi/all-with-students');
        setRecords(res.data?.data || []);
      } else if (isSiswa) {
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
      name_lengkap: user?.username || '',
      kelas: '',
      jurusan: '',
      hari: ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'][new Date().getDay()],
      tanggal: new Date().toISOString().split('T')[0],
      status: 'accepted',
      keterangan: 'hadir',
      keterangan_tidak_hadir: '',
      keterangan_dispen: '',
      file_dispen: ''
    });
    setError('');
    setIsModalOpen(true);
  };

  const openEditModal = (record: AbsensiRecord) => {
    setEditingId(record.Id || record.id || null);
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
    setError('');
    setIsModalOpen(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormLoading(true);
    setError('');
    try {
      const payload = {
        ...formData,
        student_id: studentId || '',
        tanggal: new Date(formData.tanggal).toISOString(),
      };

      if (editingId) {
        await api.patch(`/absensi/${editingId}`, payload);
      } else {
        if (!studentId && isSiswa) {
          setError('Data siswa belum tersedia. Pastikan kamu sudah terdaftar sebagai siswa.');
          setFormLoading(false);
          return;
        }
        await api.post('/absensi', payload);
      }
      setIsModalOpen(false);
      fetchAttendance();
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data?.message || 'Failed to save attendance.');
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
    const s = status?.toLowerCase();
    if (s === 'accepted' || s === 'h') {
      return <span className="flex w-fit items-center gap-2 px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"><UserCheck size={14}/> Hadir</span>;
    }
    if (s === 'permissions' || s === 'i') {
      return <span className="flex w-fit items-center gap-2 px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest bg-indigo-500/10 text-indigo-400 border border-indigo-500/20"><Calendar size={14}/> Izin</span>;
    }
    return <span className="flex w-fit items-center gap-2 px-3 py-1.5 rounded-full text-[10px] font-black uppercase tracking-widest bg-red-500/10 text-red-400 border border-red-500/20"><UserX size={14}/> Absen</span>;
  };

  return (
    <div className="space-y-8 animate-fade-in-up">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-6">
        <div className="space-y-2">
          <h1 className="text-2xl font-black text-white leading-tight flex items-center gap-3 uppercase tracking-tight">
            <Calendar size={24} className="text-indigo-400" /> Presensi Siswa
          </h1>
          <p className="text-slate-400 text-sm font-medium leading-relaxed">Kelola dan lihat riwayat absensi harian kelas.</p>
        </div>
        
        {isSiswa && (
          <Button 
            onClick={openAddModal} 
            icon={<Plus size={18} />} 
            className="font-black tracking-wide"
          >
            Record Attendance
          </Button>
        )}
      </div>

      {/* Table Section */}
      <div className="glass rounded-3xl border border-white/5 shadow-2xl overflow-hidden">
        <div className="px-8 py-6 border-b border-white/5 bg-white/2 flex items-center justify-between">
          <h3 className="text-sm font-black text-white uppercase tracking-widest flex items-center gap-2">
            <Filter size={16} className="text-slate-500" /> Log Kehadiran
          </h3>
          <div className="relative hidden sm:block">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" size={14} />
            <input type="text" placeholder="Cari nama siswa..." className="bg-white/5 border border-white/10 rounded-full py-1.5 pl-9 pr-4 text-xs text-white focus:outline-none focus:ring-1 focus:ring-indigo-500/50 w-48 transition-all" />
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-white/1">
                <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em]">Date & Day</th>
                <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em]">Student Name</th>
                <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em]">Status</th>
                <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em] hidden md:table-cell">Keterangan</th>
                <th className="px-8 py-5 text-[10px] font-black text-slate-500 uppercase tracking-[0.2em] text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {loading ? (
                <tr>
                  <td colSpan={5} className="px-8 py-20 text-center">
                    <div className="animate-spin w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full mx-auto mb-4" />
                    <p className="text-xs text-slate-500 font-bold uppercase tracking-widest">Loading records...</p>
                  </td>
                </tr>
              ) : records.length === 0 ? (
                <tr>
                  <td colSpan={5} className="px-8 py-20 text-center space-y-4">
                    <div className="w-16 h-16 bg-white/5 rounded-2xl flex items-center justify-center mx-auto text-slate-600">
                      <Calendar size={32} />
                    </div>
                    <p className="font-black text-lg text-slate-400 uppercase tracking-tight">No attendance records found.</p>
                    <p className="text-xs text-slate-500 font-medium">Silakan lakukan absensi terlebih dahulu.</p>
                  </td>
                </tr>
              ) : (
                records.map((record, i) => {
                  const id = record.Id || record.id || i.toString();
                  return (
                    <tr key={id} className="hover:bg-white/2 group transition-all">
                      <td className="px-8 py-6">
                        <div className="flex flex-col">
                          <span className="text-sm font-black text-white">
                            {record.Tanggal ? new Date(record.Tanggal).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : 'No date'}
                          </span>
                          <span className="text-[10px] text-slate-500 font-black uppercase tracking-widest leading-none mt-1">
                            <Clock size={12} className="inline mr-1 text-slate-700" /> {record.Hari || record.hari || '—'}
                          </span>
                        </div>
                      </td>
                      <td className="px-8 py-6">
                        <div className="flex items-center gap-3">
                          <div className="w-10 h-10 rounded-xl bg-linear-to-br from-indigo-500/20 to-purple-500/20 border border-white/10 flex items-center justify-center text-indigo-300 font-black text-sm">
                            {(record.NameLengkap || record.name_lengkap || record.StudentFullName)?.charAt(0) || <User size={16} />}
                          </div>
                          <span className="text-sm font-black text-white group-hover:text-indigo-400 transition-colors">
                            {record.NameLengkap || record.name_lengkap || record.StudentFullName || 'Unknown Student'}
                          </span>
                        </div>
                      </td>
                      <td className="px-8 py-6">
                        {getStatusBadge(record.Status || record.status || '')}
                      </td>
                      <td className="px-8 py-6 hidden md:table-cell text-slate-400 text-xs font-medium italic leading-relaxed max-w-xs truncate">
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
                  )
                })
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
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="p-4 bg-red-500/10 border border-red-500/20 text-red-500 rounded-2xl text-sm font-semibold leading-normal flex items-center gap-3 animate-head-shake">
              <span className="text-lg">⚠</span> {error}
            </div>
          )}

          <div className="space-y-4">
            <Input 
              label="Nama Lengkap" 
              value={formData.name_lengkap} 
              onChange={e => setFormData({...formData, name_lengkap: e.target.value})} 
              required 
              placeholder="Masukkan nama lengkap kamu"
            />
            
            <div className="grid grid-cols-2 gap-6">
              <Input 
                label="Kelas" 
                value={formData.kelas} 
                onChange={e => setFormData({...formData, kelas: e.target.value})} 
                required 
                placeholder="Contoh: X"
              />
              <Input 
                label="Jurusan" 
                value={formData.jurusan} 
                onChange={e => setFormData({...formData, jurusan: e.target.value})} 
                required 
                placeholder="Contoh: RPL"
              />
            </div>

            <div className="grid grid-cols-2 gap-6">
              <Input 
                label="Hari" 
                as="select"
                value={formData.hari} 
                onChange={e => setFormData({...formData, hari: e.target.value})} 
                required 
              >
                {['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu', 'Minggu'].map(h => (
                  <option key={h} value={h}>{h}</option>
                ))}
              </Input>
              <Input 
                label="Tanggal" 
                type="date"
                value={formData.tanggal} 
                onChange={e => setFormData({...formData, tanggal: e.target.value})} 
                required 
              />
            </div>

            <Input
              label="Keterangan Dasar"
              as="select"
              value={formData.keterangan}
              onChange={e => setFormData({...formData, keterangan: e.target.value})}
              required
            >
              <option value="hadir">Hadir</option>
              <option value="tidak hadir">Tidak Hadir (Sakit/Alpa)</option>
              <option value="izin">Izin</option>
              <option value="dispen">Dispensasi</option>
            </Input>

            {formData.keterangan !== 'hadir' && (
              <div className="space-y-4 pt-4 border-t border-white/5 mt-4 animate-fade-in-up">
                <Input 
                  label="Keterangan Tidak Hadir / Alasan" 
                  value={formData.keterangan_tidak_hadir} 
                  onChange={e => setFormData({...formData, keterangan_tidak_hadir: e.target.value})} 
                  placeholder="e.g. Sakit perut / Keperluan keluarga"
                />
                <div className="grid grid-cols-2 gap-6">
                  <Input 
                    label="Keterangan Dispen" 
                    value={formData.keterangan_dispen} 
                    onChange={e => setFormData({...formData, keterangan_dispen: e.target.value})} 
                    placeholder="e.g. Lomba OSN"
                  />
                  <Input 
                    label="File Dispen (Nama File)" 
                    value={formData.file_dispen} 
                    onChange={e => setFormData({...formData, file_dispen: e.target.value})} 
                    placeholder="upload_file.jpg"
                  />
                </div>
              </div>
            )}
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
              {editingId ? 'Simpan Perubahan' : 'Save Record'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
