export interface Group {
  id: number;
  name: string;
  services: Service[];
}

export interface Service {
  id: number;
  name: string;
  beans: Promise<BeanStatus>;
  graph: Promise<Graph>;
}

export interface BeanStatus {
  current: Bean;
  beans: Bean[];
}

export interface Bean {
  state: number;
  time: number;
}

export interface Graph {
  x: number;
  y: number;
}
