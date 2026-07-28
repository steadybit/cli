// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

export type FetchAdviceRequest = {
  environmentName: string;
  query?: string;
  offset?: number;
};

export type FetchAdviceResponse = {
  totalItems: number;
  nextOffset?: number;
  items: AdviceStatus[];
};

export type AdviceStatus = Record<string, any> & {
  target: Record<string, any> & {
    type: string;
    reference: string;
    label: string;
  };
  advice: Record<string, any> & {
    type: string;
    label: string;
    status: string;
  };
  url: string;
};
