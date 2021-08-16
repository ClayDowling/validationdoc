using System;
using System.Runtime.CompilerServices;

namespace gameoflife
{
    public class Board
    {
        private Cell[,] _cells;

        /// <summary>
        /// Initializes the board to a given size.
        /// </summary>
        /// <param name="x"></param>
        /// <param name="y"></param>
        /// <param name="initializer">
        /// Function which should return true if the cell starts live, false if the cell starts dead.
        /// </param>
        public Board(int x, int y, Func<bool> initializer)
        {
            MaxX = x;
            MaxY = y;
            _cells = new Cell[x, y];

            for (int i = 0; i < x; i++)
            {
                for (int j = 0; j < y; j++)
                {
                    _cells[i, j] = new Cell(initializer());
                }
            }
        }
        
        public int MaxX { get; }
        public int MaxY { get; }

        public Cell this[int x, int y]
        {
            get => _cells[x, y];
        }
    
    }
}