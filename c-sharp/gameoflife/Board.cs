using System.Runtime.CompilerServices;

namespace gameoflife
{
    public class Board
    {
        private Cell[,] _cells;

        public Board(int x, int y)
        {
            MaxX = x;
            MaxY = y;
            _cells = new Cell[x, y];
        }
        
        public int MaxX { get; }
        public int MaxY { get; }

        public Cell this[int x, int y]
        {
            get => _cells[x, y];
        }
    
    }
}