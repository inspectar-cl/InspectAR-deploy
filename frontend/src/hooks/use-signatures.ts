import { useEffect, useState } from 'react';
import type { SavedSignature } from '@/components/dashboard/reportes-tecnicos/SignatureDialog';

const STORAGE_KEY = 'inspectar.signatures.v1';

export function useSignatures() {
  const [signatures, setSignatures] = useState<SavedSignature[]>([]);

  useEffect(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) setSignatures(JSON.parse(raw));
    } catch {
      // noop
    }
  }, []);

  const saveAll = (items: SavedSignature[]) => {
    setSignatures(items);
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
    } catch {
      // noop
    }
  };

  const addSignature = (sig: SavedSignature) => {
    const items = [sig, ...signatures].slice(0, 25); 
    saveAll(items);
  };

  const removeSignature = (id: string) => {
    saveAll(signatures.filter(s => s.id !== id));
  };

  return { signatures, addSignature, removeSignature };
}
