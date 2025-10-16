import { useEffect, useState } from 'react';
import type { SavedSignature } from '@/components/dashboard/reportes-tecnicos/SignatureDialog';

const STORAGE_KEY = 'inspectar.signatures.v1';

export function useSignatures(): {
    signatures: SavedSignature[];
    addSignature: (sig: SavedSignature) => void;
    removeSignature: (id: string) => void;
} {
  const [signatures, setSignatures] = useState<SavedSignature[]>([]);

  useEffect(():void => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) {
        const parsedSignatures = JSON.parse(raw) as SavedSignature[];
        setSignatures(parsedSignatures);
      }
    } catch {
      // noop
    }
  }, []);

  const saveAll = (items: SavedSignature[]):void => {
    setSignatures(items);
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(items));
    } catch {
      // noop
    }
  };

  const addSignature = (sig: SavedSignature): void => {
    const items = [sig, ...signatures].slice(0, 25); 
    saveAll(items);
  };

  const removeSignature = (id: string): void => {
    saveAll(signatures.filter(s => s.id !== id));
  };

  return { signatures, addSignature, removeSignature };
}
