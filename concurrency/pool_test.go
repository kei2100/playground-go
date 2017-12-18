package concurrency

// sync.Poolは一時的なオブジェクトのPool
// groutineセーフにオブジェクトを共有するPoolからGetしたり、Putしたりすることができる。
// 内部は弱参照になっており、sync.Pool内でしか参照されていないオブジェクトは、
// 予告なしに解放される可能性がある。
